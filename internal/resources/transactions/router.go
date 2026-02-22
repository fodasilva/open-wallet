package transactions

import (
	"log"

	"github.com/felipe1496/open-wallet/db"

	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	db, err := db.Conn(utils.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	jwtService := services.NewJWTService()
	handler := NewHandler(db)
	transactionsGroup := router.Group("/api/v1/transactions")
	{
		transactionsGroup.GET("/entries",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryOptsMiddleware(),
			handler.ListEntries)
		transactionsGroup.DELETE("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.DeleteTransaction)
		transactionsGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.CreateTransaction)
		transactionsGroup.PATCH("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.UpdateTransaction)
	}
}
