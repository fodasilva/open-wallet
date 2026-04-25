package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/handlers"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func SetupTransactionsRoutes(r *gin.Engine, f *factory.Factory, redisClient *redis.Client, cfg *infra.Config) {
	jwtService := f.JWTService()
	transactionsHandler := handlers.NewHandler(f.TransactionsUseCases())
	transactionsFilterConfig := querybuilder.ParseConfig{
		AllowedFields: map[string]querybuilder.FieldConfig{
			"category_id":    {AllowedOperators: []string{"eq", "in"}},
			"type":           {AllowedOperators: []string{"eq", "in"}},
			"reference_date": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
			"amount":         {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
			"id":             {AllowedOperators: []string{"eq", "in"}},
			"user_id":        {AllowedOperators: []string{"eq", "in"}},
			"period":         {AllowedOperators: []string{"eq", "in", "gte", "lte"}},
			"created_at":     {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		},
		AllowedSortFields: []string{"reference_date", "amount", "id", "created_at"},
	}

	summaryFilterConfig := querybuilder.ParseConfig{
		AllowedFields: map[string]querybuilder.FieldConfig{
			"period":        {AllowedOperators: []string{"eq", "in", "gte", "lte"}},
			"total_expense": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
			"total_income":  {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
			"total_balance": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		},
		AllowedSortFields: []string{"period", "total_expense", "total_income", "total_balance"},
	}

	transactionsGroup := r.Group("/api/v1/transactions")
	txMax, txWin := cfg.RateLimits.XS()
	{
		transactionsGroup.GET("/entries",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(transactionsFilterConfig),
			transactionsHandler.ListEntries)
		transactionsGroup.GET("/summary",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(summaryFilterConfig),
			transactionsHandler.Summary)
		transactionsGroup.DELETE("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			transactionsHandler.DeleteTransaction)
		transactionsGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(redisClient, txMax, txWin, "transactions:create"),
			transactionsHandler.CreateTransaction)
		transactionsGroup.PATCH("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(redisClient, txMax, txWin, "transactions:update"),
			transactionsHandler.UpdateTransaction)
	}
}
