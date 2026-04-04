package routes

import (
	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRecurrencesRoutes(r *gin.Engine, f *factory.Factory, redisClient *redis.Client, cfg *infra.Config) {
	jwtService := f.JWTService()
	recurrencesHandler := recurrences.NewHandler(f.RecurrencesUseCases())
	recurrencesGroup := r.Group("/api/v1/recurrences")
	{
		recurrencesGroup.GET("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryOptsMiddleware(),
			recurrencesHandler.List)
		recurrencesGroup.DELETE("/:id",
			middlewares.RequireAuthMiddleware(jwtService),
			recurrencesHandler.DeleteByID)
		recurrencesGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "POST:/api/v1/recurrences"),
			recurrencesHandler.Create)
		recurrencesGroup.PATCH("/:id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "PATCH:/api/v1/recurrences/:id"),
			recurrencesHandler.Update)
		recurrencesGroup.POST("/:period",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 20, 60000, "POST:/api/v1/recurrences/:period"),
			recurrencesHandler.Prepare)
	}
}
