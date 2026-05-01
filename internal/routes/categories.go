package routes

import (
	"net/http"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/categories/handlers"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func SetupCategoriesRoutes(mux *http.ServeMux, f *factory.Factory, cfg *infra.Config) {
	jwtService := f.JWTService()
	categoriesHandler := handlers.NewHandler(f.CategoriesUseCases())
	catMax, catWin := cfg.RateLimits.XS()

	mux.Handle("POST /api/v1/categories", httputil.Chain(
		categoriesHandler.Create,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.NewRateLimitMiddleware(f.CacheService(), catMax, catWin, "categories:create"),
	))
	mux.Handle("GET /api/v1/categories", httputil.Chain(
		categoriesHandler.List,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.QueryBuilderMiddleware(handlers.CategoriesFilterConfig),
	))
	mux.Handle("DELETE /api/v1/categories/{category_id}", httputil.Chain(
		categoriesHandler.DeleteByID,
		middlewares.RequireAuthMiddleware(jwtService),
	))
	mux.Handle("GET /api/v1/categories/{period}", httputil.Chain(
		categoriesHandler.ListCategoryAmountPerPeriod,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.QueryBuilderMiddleware(handlers.PeriodCategoriesFilterConfig),
	))
	mux.Handle("PATCH /api/v1/categories/{category_id}", httputil.Chain(
		categoriesHandler.Update,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.NewRateLimitMiddleware(f.CacheService(), catMax, catWin, "categories:update"),
	))
}
