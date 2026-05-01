package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/categories/handlers"
)

func SetupCategoriesRoutes(r *gin.Engine, f *factory.Factory, cfg *infra.Config) {
	jwtService := f.JWTService()
	categoriesHandler := handlers.NewHandler(f.CategoriesUseCases())
	categoriesGroup := r.Group("/api/v1/categories")
	catMax, catWin := cfg.RateLimits.XS()
	{
		categoriesGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(f.CacheService(), catMax, catWin, "categories:create"),
			categoriesHandler.Create)
		categoriesGroup.GET("", middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(handlers.CategoriesFilterConfig),
			categoriesHandler.List)
		categoriesGroup.DELETE("/:category_id",
			middlewares.RequireAuthMiddleware(jwtService),
			categoriesHandler.DeleteByID)
		categoriesGroup.GET("/:period",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(handlers.PeriodCategoriesFilterConfig),
			categoriesHandler.ListCategoryAmountPerPeriod)
		categoriesGroup.PATCH("/:category_id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(f.CacheService(), catMax, catWin, "categories:update"),
			categoriesHandler.Update)
	}
}
