package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
)

func SetupRoutes(r *gin.Engine, f *factory.Factory, cfg *infra.Config) {
	SetupAuthRoutes(r, f, cfg)
	SetupCategoriesRoutes(r, f, cfg)
	SetupTransactionsRoutes(r, f, cfg)
	SetupRecurrencesRoutes(r, f, cfg)
}
