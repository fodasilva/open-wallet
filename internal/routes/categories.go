package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/categories/handlers"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func SetupCategoriesRoutes(r *gin.Engine, f *factory.Factory, redisClient *redis.Client, cfg *infra.Config) {
	jwtService := f.JWTService()
	categoriesHandler := handlers.NewHandler(f.CategoriesUseCases())
	categoriesFilterConfig := &querybuilder.ParseConfig{
		AllowedFields: map[string]querybuilder.FieldConfig{
			"name":       {AllowedOperators: []string{"eq", "like"}},
			"color":      {AllowedOperators: []string{"eq"}},
			"created_at": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
			"id":         {AllowedOperators: []string{"eq"}},
			"user_id":    {AllowedOperators: []string{"eq"}},
		},
		AllowedSortFields: []string{"name", "created_at", "id"},
	}

	periodCategoriesFilterConfig := &querybuilder.ParseConfig{
		AllowedFields: map[string]querybuilder.FieldConfig{
			"name":         {AllowedOperators: []string{"eq", "like"}},
			"color":        {AllowedOperators: []string{"eq"}},
			"total_amount": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
			"period":       {AllowedOperators: []string{"eq"}},
			"id":           {AllowedOperators: []string{"eq"}},
			"user_id":      {AllowedOperators: []string{"eq"}},
		},
		AllowedSortFields: []string{"name", "total_amount", "period", "id"},
	}

	categoriesGroup := r.Group("/api/v1/categories")
	{
		categoriesGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "POST:/api/v1/categories"),
			categoriesHandler.Create)
		categoriesGroup.GET("", middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(categoriesFilterConfig),
			categoriesHandler.List)
		categoriesGroup.DELETE("/:category_id",
			middlewares.RequireAuthMiddleware(jwtService),
			categoriesHandler.DeleteByID)
		categoriesGroup.GET("/:period",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(periodCategoriesFilterConfig),
			categoriesHandler.ListCategoryAmountPerPeriod)
		categoriesGroup.PATCH("/:category_id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "PATCH:/api/v1/categories/:category_id"),
			categoriesHandler.Update)
	}
}
