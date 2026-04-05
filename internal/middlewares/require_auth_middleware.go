package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func RequireAuthMiddleware(JWTService services.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			apiErr := utils.NewHTTPError(http.StatusUnauthorized, "missing token")
			ctx.JSON(apiErr.StatusCode, apiErr)
			ctx.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := JWTService.ValidateToken(tokenString)
		if err != nil {
			apiErr := err.(*utils.HTTPError)
			ctx.JSON(apiErr.StatusCode, apiErr)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
