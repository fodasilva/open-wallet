package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	docs "github.com/felipe1496/open-wallet/docs"

	"github.com/felipe1496/open-wallet/internal/resources/auth"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func DelayMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		delayStr := os.Getenv("DELAY")
		if delayStr != "" {
			delayNum, err := strconv.Atoi(delayStr)
			if err == nil {
				time.Sleep(time.Duration(delayNum) * time.Millisecond)
			}
		}
		c.Next()
	}
}

// @title Money API
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	// add swagger
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(DelayMiddleware())

	err := godotenv.Load()
	utils.LoadEnvs()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	origins := os.Getenv("ORIGINS")
	if origins == "" {
		log.Fatal("ORIGINS cannot be empty")
	} else {
		log.Println("ORIGINS:", origins)
	}

	originsList := strings.Split(origins, ",")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     originsList,
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth.Router(r)
	transactions.Router(r)
	categories.Router(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
