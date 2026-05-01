package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return strings.Split(ip, ":")[0]
}

func NewRateLimitMiddleware(cache services.CacheService, maxRequests int, windowMilliseconds int, prefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)
			key := fmt.Sprintf("%s:%s", prefix, ip)

			currentCount, err := cache.Incr(r.Context(), "rate_limit", key, time.Duration(windowMilliseconds)*time.Millisecond)
			if err != nil {
				log.Printf("rate limit error: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}

			if currentCount > maxRequests {
				apiErr := util.NewHTTPError(http.StatusTooManyRequests, "You have exceeded the rate limit")
				httputil.JSON(w, apiErr.StatusCode, apiErr)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
