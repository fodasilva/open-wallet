package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/auth/handlers"
)

func SetupAuthRoutes(r *gin.Engine, f *factory.Factory, cfg *infra.Config) {
	authHandler := handlers.NewHandler(f.AuthUseCases(), f.JWTService())
	authGroup := r.Group("/api/v1/auth")
	authMax, authWin := cfg.RateLimits.XS()
	{
		authGroup.POST("/login/google",
			middlewares.NewRateLimitMiddleware(f.CacheService(), authMax, authWin, "auth:google-login"),
			authHandler.CreateLoginWithGoogle)
	}
}
