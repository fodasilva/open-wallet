package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/handlers"
)

func SetupRecurrencesRoutes(r *gin.Engine, f *factory.Factory, cfg *infra.Config) {
	jwtService := f.JWTService()
	recurrencesHandler := handlers.NewHandler(f.RecurrencesUseCases())

	recurrencesGroup := r.Group("/api/v1/recurrences")
	recMax, recWin := cfg.RateLimits.XS()
	prepMax, prepWin := cfg.RateLimits.SM()
	{
		recurrencesGroup.GET("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryBuilderMiddleware(handlers.RecurrencesFilterConfig),
			recurrencesHandler.List)
		recurrencesGroup.DELETE("/:id",
			middlewares.RequireAuthMiddleware(jwtService),
			recurrencesHandler.Delete)
		recurrencesGroup.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(f.CacheService(), recMax, recWin, "recurrences:create"),
			recurrencesHandler.Create)
		recurrencesGroup.PATCH("/:id",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(f.CacheService(), recMax, recWin, "recurrences:update"),
			recurrencesHandler.Update)
		recurrencesGroup.POST("/:period",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.NewRateLimitMiddleware(f.CacheService(), prepMax, prepWin, "recurrences:prepare"),
			recurrencesHandler.Prepare)
	}
}
