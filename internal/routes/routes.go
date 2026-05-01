package routes

import (
	"net/http"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
)

func SetupRoutes(mux *http.ServeMux, f *factory.Factory, cfg *infra.Config) {
	SetupAuthRoutes(mux, f, cfg)
	SetupCategoriesRoutes(mux, f, cfg)
	SetupTransactionsRoutes(mux, f, cfg)
	SetupRecurrencesRoutes(mux, f, cfg)
}
