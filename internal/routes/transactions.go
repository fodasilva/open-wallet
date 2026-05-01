package routes

import (
	"net/http"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/handlers"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func SetupTransactionsRoutes(mux *http.ServeMux, f *factory.Factory, cfg *infra.Config) {
	jwtService := f.JWTService()
	transactionsHandler := handlers.NewHandler(f.TransactionsUseCases())
	txMax, txWin := cfg.RateLimits.XS()

	mux.Handle("GET /api/v1/transactions/entries", httputil.Chain(
		transactionsHandler.ListEntries,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.QueryBuilderMiddleware(handlers.TransactionsFilterConfig),
	))
	mux.Handle("GET /api/v1/transactions/summary", httputil.Chain(
		transactionsHandler.Summary,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.QueryBuilderMiddleware(handlers.SummaryFilterConfig),
	))
	mux.Handle("DELETE /api/v1/transactions/{transaction_id}", httputil.Chain(
		transactionsHandler.DeleteTransaction,
		middlewares.RequireAuthMiddleware(jwtService),
	))
	mux.Handle("POST /api/v1/transactions", httputil.Chain(
		transactionsHandler.CreateTransaction,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.NewRateLimitMiddleware(f.CacheService(), txMax, txWin, "transactions:create"),
	))
	mux.Handle("PATCH /api/v1/transactions/{transaction_id}", httputil.Chain(
		transactionsHandler.UpdateTransaction,
		middlewares.RequireAuthMiddleware(jwtService),
		middlewares.NewRateLimitMiddleware(f.CacheService(), txMax, txWin, "transactions:update"),
	))
}
