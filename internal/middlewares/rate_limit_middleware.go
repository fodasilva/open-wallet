package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func rateLimit(redisClient *redis.Client, maxRequests int, windowMilliseconds int, prefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		key := fmt.Sprintf("rate_limit:%s:%s", prefix, ip)

		val, err := redisClient.Incr(context.Background(), key).Result()
		if err != nil {
			ctx.Next()
			return
		}

		if val == 1 {
			redisClient.Expire(context.Background(), key, time.Duration(windowMilliseconds)*time.Millisecond)
		}

		if val > int64(maxRequests) {
			apiErr := utils.NewHTTPError(http.StatusTooManyRequests, "You have exceeded the rate limit")
			ctx.JSON(apiErr.StatusCode, apiErr)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func GlobalRateLimitMiddleware(redisClient *redis.Client, cfg *infra.Config) gin.HandlerFunc {
	limit := cfg.RateLimitMaxRequests
	window := time.Duration(cfg.RateLimitWindowMs) * time.Millisecond
	return rateLimit(redisClient, limit, int(window.Milliseconds()), "global")
}

func RouteRateLimitMiddleware(redisClient *redis.Client, maxRequests int, windowMilliseconds int, prefix string) gin.HandlerFunc {
	return rateLimit(redisClient, maxRequests, windowMilliseconds, prefix)
}
