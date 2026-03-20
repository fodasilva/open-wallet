package infra

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment          string
	Delay                int
	Origins              []string
	GcpProjectID         *string
	Port                 int
	DatabaseURL          string
	GoogleClientID       string
	GoogleSecret         string
	LoginRedirectURI     string
	JWTSecret            string
	RateLimitDBURL       string
	RateLimitMaxRequests int
	RateLimitWindowMs    int
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file", err)
	}

	cfg := &Config{}
	var errs []string

	steps := []func(){
		func() {
			env := os.Getenv("ENVIRONMENT")
			if env == "" {
				cfg.Environment = "dev"
				return
			}
			allowed := []string{"dev", "prod", "test"}
			if !slices.Contains(allowed, env) {
				errs = append(errs, "ENVIRONMENT must be one of: dev, prod, test")
			}
			cfg.Environment = env
		},
		func() {
			val := os.Getenv("DELAY")
			if val == "" {
				cfg.Delay = 0
				return
			}
			intVal, err := strconv.Atoi(val)
			if err != nil {
				errs = append(errs, "DELAY must be a number")
				cfg.Delay = 0
				return
			}
			cfg.Delay = intVal
		},
		func() {
			val := os.Getenv("ORIGINS")
			if val != "" {
				parts := strings.Split(val, ",")
				for _, part := range parts {
					if !utils.IsValidURL(part) {
						errs = append(errs, fmt.Sprintf("ORIGINS item %s must be a valid url", part))
					}
				}
				cfg.Origins = parts
			}
		},
		func() {
			val := os.Getenv("GCP_PROJECT_ID")
			if val != "" {
				cfg.GcpProjectID = &val
			}
		},
		func() {
			val := os.Getenv("PORT")
			if val == "" {
				cfg.Port = 8080
				return
			}
			intVal, err := strconv.Atoi(val)
			if err != nil {
				errs = append(errs, "PORT must be a number")
				cfg.Port = 8080
				return
			}
			cfg.Port = intVal
		},
		func() {
			val := os.Getenv("DATABASE_URL")
			if val == "" {
				errs = append(errs, "DATABASE_URL is required")
			}
			cfg.DatabaseURL = val
		},
		func() {
			val := os.Getenv("GOOGLE_CLIENT_ID")
			if val == "" {
				errs = append(errs, "GOOGLE_CLIENT_ID is required")
			}
			cfg.GoogleClientID = val
		},
		func() {
			val := os.Getenv("GOOGLE_CLIENT_SECRET")
			if val == "" {
				errs = append(errs, "GOOGLE_CLIENT_SECRET is required")
			}
			cfg.GoogleSecret = val
		},
		func() {
			val := os.Getenv("LOGIN_REDIRECT_URI")
			if val == "" {
				errs = append(errs, "LOGIN_REDIRECT_URI is required")
			}
			cfg.LoginRedirectURI = val
		},
		func() {
			val := os.Getenv("JWT_SECRET")
			if val == "" {
				errs = append(errs, "JWT_SECRET is required")
			}
			cfg.JWTSecret = val
		},
		func() {
			val := os.Getenv("RATE_LIMIT_DB_URL")
			if val == "" {
				errs = append(errs, "RATE_LIMIT_DB_URL is required")
			}
			cfg.RateLimitDBURL = val
		},
		func() {
			val := os.Getenv("RATE_LIMIT_MAX_REQUESTS")
			if val == "" {
				cfg.RateLimitMaxRequests = 100
				return
			}
			intVal, err := strconv.Atoi(val)
			if err != nil {
				errs = append(errs, "RATE_LIMIT_MAX_REQUESTS must be a number")
				cfg.RateLimitMaxRequests = 100
				return
			}
			cfg.RateLimitMaxRequests = intVal
		},
		func() {
			val := os.Getenv("RATE_LIMIT_WINDOW_MS")
			if val == "" {
				cfg.RateLimitWindowMs = 60000
				return
			}
			intVal, err := strconv.Atoi(val)
			if err != nil {
				errs = append(errs, "RATE_LIMIT_WINDOW_MS must be a number")
				cfg.RateLimitWindowMs = 60000
				return
			}
			cfg.RateLimitWindowMs = intVal
		},
	}

	for _, step := range steps {
		step()
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("invalid configuration: %s", strings.Join(errs, "; "))
	}

	return cfg, nil
}
