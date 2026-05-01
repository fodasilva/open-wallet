package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/handlers"
)

func SetupTransactionsRoutes(r *gin.Engine, f *factory.Factory, cfg *infra.Config) {
	jwtService := f.JWTService()
	transactionsHandler := handlers.NewHandler(f.TransactionsUseCases())
	transactionsGroup := r.Group("/api/v1/transactions")
	txMax, txWin := cfg.RateLimits.XS()
	{
		transactionsGroup.GET("/entries",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(handlers.TransactionsFilterConfig),
			transactionsHandler.ListEntries)
		transactionsGroup.GET("/summary",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(handlers.SummaryFilterConfig),
			transactionsHandler.Summary)
		transactionsGroup.DELETE("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			transactionsHandler.DeleteTransaction)
		transactionsGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(f.CacheService(), txMax, txWin, "transactions:create"),
			transactionsHandler.CreateTransaction)
		transactionsGroup.PATCH("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(f.CacheService(), txMax, txWin, "transactions:update"),
			transactionsHandler.UpdateTransaction)
	}
}
