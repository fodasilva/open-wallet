package routes

import (
	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupAuthRoutes(r *gin.Engine, f *factory.Factory, redisClient *redis.Client, cfg *infra.Config) {
	authHandler := auth.NewHandler(f.GoogleService(), f.JWTService(), f.UsersUseCases(), f.AuthUseCases())
	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/login/google",
			middlewares.RouteRateLimitMiddleware(redisClient, 5, 300000, "POST:/api/v1/auth/login/google"),
			authHandler.LoginGoogle)
	}
}
