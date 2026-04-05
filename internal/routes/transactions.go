package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
)

func SetupTransactionsRoutes(r *gin.Engine, f *factory.Factory, redisClient *redis.Client, cfg *infra.Config) {
	jwtService := f.JWTService()
	transactionsHandler := transactions.NewHandler(f.TransactionsUseCases())
	transactionsGroup := r.Group("/api/v1/transactions")
	{
		transactionsGroup.GET("/entries",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(nil),
			transactionsHandler.ListEntries)
		transactionsGroup.DELETE("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			transactionsHandler.DeleteTransaction)
		transactionsGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "POST:/api/v1/transactions"),
			transactionsHandler.CreateTransaction)
		transactionsGroup.PATCH("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 60000, "PATCH:/api/v1/transactions/:transaction_id"),
			transactionsHandler.UpdateTransaction)
	}
}
