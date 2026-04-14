package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/handlers"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func SetupRecurrencesRoutes(r *gin.Engine, f *factory.Factory, redisClient *redis.Client, cfg *infra.Config) {
	jwtService := f.JWTService()
	recurrencesHandler := handlers.NewHandler(f.RecurrencesUseCases())
	recurrencesFilterConfig := querybuilder.ParseConfig{
		AllowedFields: map[string]querybuilder.FieldConfig{
			"id":          {AllowedOperators: []string{"eq", "in"}},
			"category_id": {AllowedOperators: []string{"eq", "in"}},
			"name":        {AllowedOperators: []string{"eq", "like", "in"}},
			"user_id":     {AllowedOperators: []string{"eq", "in"}},
			"created_at":  {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
			"amount":      {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		},
		AllowedSortFields: []string{"name", "created_at", "id"},
	}

	recurrencesGroup := r.Group("/api/v1/recurrences")
	recMax, recWin := cfg.RateLimits.XS()
	prepMax, prepWin := cfg.RateLimits.SM()
	{
		recurrencesGroup.GET("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(recurrencesFilterConfig),
			recurrencesHandler.List)
		recurrencesGroup.DELETE("/:id",
			middlewares.RequireAuthMiddleware(jwtService),
			recurrencesHandler.Delete)
		recurrencesGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(redisClient, recMax, recWin, "recurrences:create"),
			recurrencesHandler.Create)
		recurrencesGroup.PATCH("/:id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(redisClient, recMax, recWin, "recurrences:update"),
			recurrencesHandler.Update)
		recurrencesGroup.POST("/:period",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(redisClient, prepMax, prepWin, "recurrences:prepare"),
			recurrencesHandler.Prepare)
	}
}
