package routes

import (
	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(r *gin.Engine, f *factory.Factory, redisClient *redis.Client, cfg *infra.Config) {
	SetupAuthRoutes(r, f, redisClient, cfg)
	SetupCategoriesRoutes(r, f, redisClient, cfg)
	SetupTransactionsRoutes(r, f, redisClient, cfg)
	SetupRecurrencesRoutes(r, f, redisClient, cfg)
}
