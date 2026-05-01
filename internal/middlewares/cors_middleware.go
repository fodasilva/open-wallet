package middlewares

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/felipe1496/open-wallet/infra"
)

func CorsMiddleware(cfg *infra.Config) func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.Origins,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
	return c.Handler
}
