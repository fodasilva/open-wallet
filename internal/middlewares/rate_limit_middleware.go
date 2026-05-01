package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func NewRateLimitMiddleware(cache services.CacheService, maxRequests int, windowMilliseconds int, prefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		key := fmt.Sprintf("%s:%s", prefix, ip)

		currentCount, err := cache.Incr(ctx.Request.Context(), "rate_limit", key, time.Duration(windowMilliseconds)*time.Millisecond)
		if err != nil {
			log.Printf("rate limit error: %v\n", err)
			ctx.Next()
			return
		}

		if currentCount > maxRequests {
			apiErr := utils.NewHTTPError(http.StatusTooManyRequests, "You have exceeded the rate limit")
			ctx.JSON(apiErr.StatusCode, apiErr)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
