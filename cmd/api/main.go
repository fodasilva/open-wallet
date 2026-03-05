package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	docs "github.com/felipe1496/open-wallet/docs"
	"github.com/felipe1496/open-wallet/trace"

	"github.com/felipe1496/open-wallet/db"

	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/auth"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Open Wallet API
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cleanupTracer := setupTracer()
	defer cleanupTracer()

	dbConn, redisClient := setupPersistence()

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(middlewares.DelayMiddleware())
	r.Use(middlewares.CorsMiddleware())
	r.Use(middlewares.TraceMiddleware("open-wallet-service"))
	r.Use(gin.Recovery())
	r.Use(middlewares.GlobalRateLimitMiddleware(redisClient))

	auth.Router(r, dbConn, redisClient)
	transactions.Router(r, dbConn, redisClient)
	categories.Router(r, dbConn, redisClient)

	r.Run(fmt.Sprintf(":%d", utils.AppConfig.Port))
}

func setupTracer() func() {
	tp, err := trace.InitTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	return func() {
		_ = tp.Shutdown(context.Background())
	}
}

func setupPersistence() (*sql.DB, *redis.Client) {
	dbConn, err := db.Conn(utils.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	opts, err := redis.ParseURL(utils.AppConfig.RateLimitDBURL)
	if err != nil {
		log.Fatalf("failed to parse redis url for rate limit: %v", err)
	}
	redisClient := redis.NewClient(opts)

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return dbConn, redisClient
}
