package recurrences

import (
	"database/sql"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/redis/go-redis/v9"

	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/services"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine, db *sql.DB, redisClient *redis.Client, cfg *infra.Config) {
	jwtService := services.NewJWTService(cfg)
	handler := NewHandler(db)
	recurrencesGroup := router.Group("/api/v1/recurrences")
	{
		recurrencesGroup.GET("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryOptsMiddleware(),
			handler.List)
		recurrencesGroup.DELETE("/:id",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.DeleteByID)
		recurrencesGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "POST:/api/v1/recurrences"),
			handler.Create)
		recurrencesGroup.PATCH("/:id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "PATCH:/api/v1/recurrences/:id"),
			handler.Update)
		recurrencesGroup.POST("/:period",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 20, 60000, "POST:/api/v1/recurrences/:period"),
			handler.Prepare)
	}
}
