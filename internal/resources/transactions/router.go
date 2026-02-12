package transactions

import (
	"log"
	"os"
	"time"

	"github.com/felipe1496/open-wallet/db"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	db, err := db.Conn(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	jwtService := services.NewJWTService()
	handler := NewHandler(db)
	transactionsGroup := router.Group("/api/v1/transactions")
	rateLimiter := middlewares.NewRateLimiter(db)
	{
		transactionsGroup.GET("/entries",
			middlewares.RequireAuthMiddleware(jwtService),
			rateLimiter.Middleware(utils.EnvConfig.GETRpmLimit, 1*time.Minute),
			middlewares.QueryOptsMiddleware(),
			handler.ListEntries)
		transactionsGroup.DELETE("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			rateLimiter.Middleware(utils.EnvConfig.DELETERpmLimit, 1*time.Minute),
			handler.DeleteTransaction)
		transactionsGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			rateLimiter.Middleware(utils.EnvConfig.POSTRpmLimit, 1*time.Minute),
			handler.CreateTransaction)
		transactionsGroup.PATCH("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			rateLimiter.Middleware(utils.EnvConfig.PATCHRpmLimit, 1*time.Minute),
			handler.UpdateTransaction)
	}
}
