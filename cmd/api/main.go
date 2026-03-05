package main

import (
	"context"
	"log"
	"strconv"
	"time"

	docs "github.com/felipe1496/open-wallet/docs"
	"github.com/felipe1496/open-wallet/trace"

	"github.com/felipe1496/open-wallet/db"

	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/auth"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/redis/go-redis/v9"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func DelayMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.AppConfig.Delay != 0 {
			time.Sleep(time.Duration(utils.AppConfig.Delay) * time.Millisecond)
		}
		c.Next()
	}
}

// @title Money API
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	tp, err := trace.InitTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	dbConn, err := db.Conn(utils.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	opts, err := redis.ParseURL(utils.AppConfig.RateLimitDBURL)
	if err != nil {
		log.Fatalf("failed to parse redis url for rate limite: %v", err)
	}
	redisClient := redis.NewClient(opts)

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	// add swagger
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(DelayMiddleware())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     utils.AppConfig.Origins,
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middlewares.Tracing("open-wallet-service"))
	r.Use(gin.Recovery())
	r.Use(middlewares.GlobalRateLimitMiddleware(redisClient))

	auth.Router(r, dbConn, redisClient)
	transactions.Router(r, dbConn, redisClient)
	categories.Router(r, dbConn, redisClient)

	port := strconv.Itoa(utils.AppConfig.Port)

	r.Run(":" + port)
}
