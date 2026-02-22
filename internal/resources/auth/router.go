package auth

import (
	"log"

	"github.com/felipe1496/open-wallet/db"

	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	db, err := db.Conn(utils.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	handler := NewHandler(db, services.NewGoogleService(), services.NewJWTService())
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login/google", handler.LoginGoogle)
	}
}
