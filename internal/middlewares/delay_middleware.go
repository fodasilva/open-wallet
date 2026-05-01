package middlewares

import (
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/infra"
)

func DelayMiddleware(cfg *infra.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Delay > 0 {
				time.Sleep(time.Duration(cfg.Delay) * time.Millisecond)
			}
			next.ServeHTTP(w, r)
		})
	}
}
