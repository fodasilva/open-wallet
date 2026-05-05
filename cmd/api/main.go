package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/felipe1496/open-wallet/docs"
	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/routes"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

// @title Open Wallet API
// @version 1.0
// @description This is the Open Wallet API.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := infra.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	cleanupTracer := setupTracer(cfg)
	defer cleanupTracer()

	dbConn, cleanupPersistence := setupPersistence(cfg)
	defer cleanupPersistence()

	f := factory.NewFactory(dbConn, cfg)

	mux := http.NewServeMux()

	if cfg.Environment == "dev" {
		mux.Handle("GET /api-docs/", httpSwagger.WrapHandler)
	}

	routes.SetupRoutes(mux, f, cfg)

	globalMax, globalWin := cfg.RateLimits.MD()

	handler := httputil.Chain(
		mux.ServeHTTP,
		middlewares.TimeoutMiddleware(time.Duration(cfg.RequestTimeout)*time.Millisecond),
		middlewares.RecoveryMiddleware(),
		middlewares.DelayMiddleware(cfg),
		middlewares.CorsMiddleware(cfg),
		middlewares.NewRateLimitMiddleware(f.CacheService(), globalMax, globalWin, "global"),
	)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Starting server on port %d", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
}

func setupPersistence(cfg *infra.Config) (*sql.DB, func()) {
	dbConn, err := infra.DBConn(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return dbConn, func() {
		_ = dbConn.Close()
	}
}

func setupTracer(cfg *infra.Config) func() {
	tp, err := infra.InitTracer(cfg)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	return func() {
		_ = tp.Shutdown(context.Background())
	}
}
