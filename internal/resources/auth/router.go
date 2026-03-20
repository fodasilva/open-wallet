package auth

import (
	"database/sql"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/redis/go-redis/v9"

	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/services"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine, db *sql.DB, redisClient *redis.Client, cfg *infra.Config) {
	googleService := services.NewGoogleService(cfg)
	jwtService := services.NewJWTService(cfg)
	handler := NewHandler(db, googleService, jwtService)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login/google",
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 300000, "POST:/api/v1/auth/login/google"),
			handler.LoginGoogle)
	}
}
