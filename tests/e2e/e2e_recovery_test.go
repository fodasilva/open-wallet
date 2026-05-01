package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func TestRecoveryMiddlewareE2E(t *testing.T) {
	mux := http.NewServeMux()

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated unhandled panic")
	})

	handler := httputil.Chain(
		panicHandler,
		middlewares.RecoveryMiddleware(),
	)

	mux.Handle("GET /panic", handler)

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response util.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "Internal server error", response.ErrorData.Message)
	assert.Equal(t, "internal server error", response.ErrorData.Type)
}
