package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/felipe1496/open-wallet/docs"
	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/routes"
)

// @title Open Wallet API
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

	dbConn, redisClient, cleanupPersistence := setupPersistence(cfg)
	defer cleanupPersistence()

	r := gin.New()
	r.Use(middlewares.DelayMiddleware(cfg))
	r.Use(middlewares.CorsMiddleware(cfg))
	r.Use(middlewares.GlobalRateLimitMiddleware(redisClient, cfg))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	f := factory.NewFactory(dbConn, cfg)
	routes.SetupRoutes(r, f, redisClient, cfg)

	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func setupPersistence(cfg *infra.Config) (*sql.DB, *redis.Client, func()) {
	dbConn, err := infra.DBConn(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	redisClient, err := infra.RedisConn(cfg.RateLimitDBURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return dbConn, redisClient, func() {
		_ = dbConn.Close()
		_ = redisClient.Close()
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
