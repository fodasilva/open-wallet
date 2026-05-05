package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func TestTimeoutMiddlewareE2E(t *testing.T) {
	tests := []struct {
		name           string
		handlerSleep   time.Duration
		timeout        time.Duration
		shouldPanic    bool
		expectedStatus int
		expectedMsg    string
		isJSONError    bool
	}{
		{
			name:           "Timeout triggered",
			handlerSleep:   200 * time.Millisecond,
			timeout:        50 * time.Millisecond,
			shouldPanic:    false,
			expectedStatus: http.StatusServiceUnavailable,
			expectedMsg:    "Service Unavailable",
			isJSONError:    false,
		},
		{
			name:           "Finishes within timeout",
			handlerSleep:   10 * time.Millisecond,
			timeout:        100 * time.Millisecond,
			shouldPanic:    false,
			expectedStatus: http.StatusOK,
			expectedMsg:    "finished",
			isJSONError:    false,
		},
		{
			name:           "Panic inside timeout",
			handlerSleep:   0,
			timeout:        100 * time.Millisecond,
			shouldPanic:    true,
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Internal server error",
			isJSONError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.shouldPanic {
					panic("simulated panic")
				}
				select {
				case <-time.After(tt.handlerSleep):
					httputil.JSON(w, http.StatusOK, map[string]string{"message": "finished"})
				case <-r.Context().Done():
					return
				}
			})

			h := httputil.Chain(handler, middlewares.RecoveryMiddleware(), middlewares.TimeoutMiddleware(tt.timeout))
			mux.Handle("GET /test", h)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.isJSONError {
				var response httputil.HTTPError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, response.ErrorData.Message)
			} else {
				assert.Contains(t, w.Body.String(), tt.expectedMsg)
			}
		})
	}
}
