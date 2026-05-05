package middlewares

import (
	"net/http"
	"time"
)

// TimeoutMiddleware returns a middleware that adds a timeout to the request context.
// It uses the standard library's http.TimeoutHandler to handle the timeout logic safely.
// If the timeout is reached, it responds with a 503 Service Unavailable status and a plain text message.
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "Service Unavailable")
	}
}
