package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func RequireAuthMiddleware(JWTService services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				apiErr := util.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
				httputil.JSON(w, apiErr.StatusCode, apiErr)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			userID, err := JWTService.ValidateToken(tokenString)
			if err != nil {
				apiErr := err.(*util.HTTPError)
				httputil.JSON(w, apiErr.StatusCode, apiErr)
				return
			}

			ctx := context.WithValue(r.Context(), util.ContextKeyUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
