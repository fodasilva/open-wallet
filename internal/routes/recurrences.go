package routes

import (
	"net/http"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/handlers"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func SetupRecurrencesRoutes(mux *http.ServeMux, f *factory.Factory, cfg *infra.Config) {
	jwtService := f.JWTService()
	recurrencesHandler := handlers.NewHandler(f.RecurrencesUseCases())
	recMax, recWin := cfg.RateLimits.XS()
	prepMax, prepWin := cfg.RateLimits.SM()

	mux.Handle("GET /api/v1/recurrences", httputil.Chain(
		recurrencesHandler.List,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.QueryBuilderMiddleware(handlers.RecurrencesFilterConfig),
	))
	mux.Handle("DELETE /api/v1/recurrences/{id}", httputil.Chain(
		recurrencesHandler.Delete,
		middlewares.RequireAuthMiddleware(jwtService),
	))
	mux.Handle("POST /api/v1/recurrences", httputil.Chain(
		recurrencesHandler.Create,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.NewRateLimitMiddleware(f.CacheService(), recMax, recWin, "recurrences:create"),
	))
	mux.Handle("PATCH /api/v1/recurrences/{id}", httputil.Chain(
		recurrencesHandler.Update,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.NewRateLimitMiddleware(f.CacheService(), recMax, recWin, "recurrences:update"),
	))
	mux.Handle("POST /api/v1/recurrences/{period}", httputil.Chain(
		recurrencesHandler.Prepare,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.NewRateLimitMiddleware(f.CacheService(), prepMax, prepWin, "recurrences:prepare"),
	))
}
