package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/felipe1496/open-wallet/docs"
	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/routes"
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

	if cfg.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	dbConn, cleanupPersistence := setupPersistence(cfg)
	defer cleanupPersistence()

	f := factory.NewFactory(dbConn, cfg)

	r := gin.New()
	r.Use(middlewares.DelayMiddleware(cfg))
	r.Use(middlewares.CorsMiddleware(cfg))
	globalMax, globalWin := cfg.RateLimits.MD()
	r.Use(middlewares.NewRateLimitMiddleware(f.CacheService(), globalMax, globalWin, "global"))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if cfg.Environment == "dev" {
		r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	routes.SetupRoutes(r, f, cfg)

	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
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
