package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/middlewares"
)

func TestRateLimitE2E(t *testing.T) {
	res := SetupTestResources(t)
	defer func() { _ = res.PostgresContainer.Terminate(context.Background()) }()
	defer func() { _ = res.RedisContainer.Terminate(context.Background()) }()
	defer func() { _ = res.DB.Close() }()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		maxRequests   int
		windowMs      int
		numRequests   int
		expectedCodes []int
	}{
		{
			name:          "Strict limit - block second request",
			maxRequests:   1,
			windowMs:      3600000,
			numRequests:   2,
			expectedCodes: []int{http.StatusOK, http.StatusTooManyRequests},
		},
		{
			name:          "Wider limit - allow multiple requests",
			maxRequests:   5,
			windowMs:      3600000,
			numRequests:   3,
			expectedCodes: []int{http.StatusOK, http.StatusOK, http.StatusOK},
		},
		{
			name:          "Limit of 0 - block everything",
			maxRequests:   0,
			windowMs:      3600000,
			numRequests:   1,
			expectedCodes: []int{http.StatusTooManyRequests},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a unique prefix/path for each test to avoid interference
			path := fmt.Sprintf("/test-limit-%d", i)
			prefix := fmt.Sprintf("test_limit_%d", i)

			r := gin.New()
			r.GET(path, middlewares.NewRateLimitMiddleware(res.RedisClient, tt.maxRequests, tt.windowMs, prefix), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			for i := 0; i < tt.numRequests; i++ {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, tt.expectedCodes[i], w.Code, "Request %d failed in test case '%s'", i+1, tt.name)

				if w.Code == http.StatusTooManyRequests {
					var resp map[string]interface{}
					require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
					errorData := resp["error"].(map[string]interface{})
					assert.Equal(t, "too many requests", errorData["type"])
				}
			}
		})
	}
}

// TestTShirtSizeIntegration verifies that our Config mapping works with the middleware
func TestRateLimitTShirtSizeIntegration(t *testing.T) {
	res := SetupTestResources(t)
	defer func() { _ = res.PostgresContainer.Terminate(context.Background()) }()
	defer func() { _ = res.RedisContainer.Terminate(context.Background()) }()

	cfg := &infra.Config{}
	cfg.RateLimits.XS = func() (int, int) { return 1, 3600000 }

	max, win := cfg.RateLimits.XS()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/xs", middlewares.NewRateLimitMiddleware(res.RedisClient, max, win, "xs"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// First request - OK
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, httptest.NewRequest(http.MethodGet, "/xs", nil))
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request - 429
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/xs", nil))
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}
