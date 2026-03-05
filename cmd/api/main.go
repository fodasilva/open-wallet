package main

import (
	"context"
	"log"
	"strconv"
	"time"

	docs "github.com/felipe1496/open-wallet/docs"
	"github.com/felipe1496/open-wallet/trace"

	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/resources/auth"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"

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

	auth.Router(r)
	transactions.Router(r)
	categories.Router(r)

	port := strconv.Itoa(utils.AppConfig.Port)

	r.Run(":" + port)
}
