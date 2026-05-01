package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func RecoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic recovered: %v\n%s", err, debug.Stack())
					httputil.JSON(w, http.StatusInternalServerError, util.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
