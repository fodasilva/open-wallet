package routes

import (
	"net/http"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/auth/handlers"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func SetupAuthRoutes(mux *http.ServeMux, f *factory.Factory, cfg *infra.Config) {
	authHandler := handlers.NewHandler(f.AuthUseCases(), f.JWTService())
	authMax, authWin := cfg.RateLimits.XS()

	mux.Handle("POST /api/v1/auth/login/google", httputil.Chain(
		authHandler.CreateLoginWithGoogle,
		middlewares.NewRateLimitMiddleware(f.CacheService(), authMax, authWin, "auth:google-login"),
	))
}
