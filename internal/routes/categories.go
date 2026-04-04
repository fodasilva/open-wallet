package routes

import (
	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupCategoriesRoutes(r *gin.Engine, f *factory.Factory, redisClient *redis.Client, cfg *infra.Config) {
	jwtService := f.JWTService()
	categoriesHandler := categories.NewHandler(f.CategoriesUseCase())
	categoriesGroup := r.Group("/api/v1/categories")
	{
		categoriesGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "POST:/api/v1/categories"),
			categoriesHandler.Create)
		categoriesGroup.GET("", middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryOptsMiddleware(),
			categoriesHandler.List)
		categoriesGroup.DELETE("/:category_id",
			middlewares.RequireAuthMiddleware(jwtService),
			categoriesHandler.DeleteByID)
		categoriesGroup.GET("/:period",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryOptsMiddleware(),
			categoriesHandler.ListCategoryAmountPerPeriod)
		categoriesGroup.PATCH("/:category_id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "PATCH:/api/v1/categories/:category_id"),
			categoriesHandler.Update)
	}
}
