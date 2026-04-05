package infra

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/felipe1496/open-wallet/internal/utils"
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

type loader struct {
	cfg  *Config
	errs []string
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file", err)
	}

	l := &loader{
		cfg: &Config{},
	}

	l.loadEnvironment()
	l.loadDelay()
	l.loadOrigins()
	l.loadGcpProjectID()
	l.loadPort()
	l.loadDatabaseURL()
	l.loadGoogleAuthConfig()
	l.loadJWTConfig()
	l.loadRateLimitConfig()

	if len(l.errs) > 0 {
		return nil, fmt.Errorf("invalid configuration: %s", strings.Join(l.errs, "; "))
	}

	return l.cfg, nil
}

func (l *loader) loadEnvironment() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		l.cfg.Environment = "dev"
		return
	}
	allowed := []string{"dev", "prod", "test"}
	if !slices.Contains(allowed, env) {
		l.errs = append(l.errs, "ENVIRONMENT must be one of: dev, prod, test")
	}
	l.cfg.Environment = env
}

func (l *loader) loadDelay() {
	val := os.Getenv("DELAY")
	if val == "" {
		l.cfg.Delay = 0
		return
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		l.errs = append(l.errs, "DELAY must be a number")
		l.cfg.Delay = 0
		return
	}
	l.cfg.Delay = intVal
}

func (l *loader) loadOrigins() {
	val := os.Getenv("ORIGINS")
	if val != "" {
		parts := strings.Split(val, ",")
		for _, part := range parts {
			if !utils.IsValidURL(part) {
				l.errs = append(l.errs, fmt.Sprintf("ORIGINS item %s must be a valid url", part))
			}
		}
		l.cfg.Origins = parts
	}
}

func (l *loader) loadGcpProjectID() {
	val := os.Getenv("GCP_PROJECT_ID")
	if val != "" {
		l.cfg.GcpProjectID = &val
	}
}

func (l *loader) loadPort() {
	val := os.Getenv("PORT")
	if val == "" {
		l.cfg.Port = 8080
		return
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		l.errs = append(l.errs, "PORT must be a number")
		l.cfg.Port = 8080
		return
	}
	l.cfg.Port = intVal
}

func (l *loader) loadDatabaseURL() {
	val := os.Getenv("DATABASE_URL")
	if val == "" {
		l.errs = append(l.errs, "DATABASE_URL is required")
	}
	l.cfg.DatabaseURL = val
}

func (l *loader) loadGoogleAuthConfig() {
	l.cfg.GoogleClientID = l.getRequired("GOOGLE_CLIENT_ID")
	l.cfg.GoogleSecret = l.getRequired("GOOGLE_CLIENT_SECRET")
	l.cfg.LoginRedirectURI = l.getRequired("LOGIN_REDIRECT_URI")
}

func (l *loader) loadJWTConfig() {
	l.cfg.JWTSecret = l.getRequired("JWT_SECRET")
}

func (l *loader) loadRateLimitConfig() {
	l.cfg.RateLimitDBURL = l.getRequired("RATE_LIMIT_DB_URL")

	maxReq := os.Getenv("RATE_LIMIT_MAX_REQUESTS")
	if maxReq == "" {
		l.cfg.RateLimitMaxRequests = 100
	} else if val, err := strconv.Atoi(maxReq); err != nil {
		l.errs = append(l.errs, "RATE_LIMIT_MAX_REQUESTS must be a number")
		l.cfg.RateLimitMaxRequests = 100
	} else {
		l.cfg.RateLimitMaxRequests = val
	}

	window := os.Getenv("RATE_LIMIT_WINDOW_MS")
	if window == "" {
		l.cfg.RateLimitWindowMs = 60000
	} else if val, err := strconv.Atoi(window); err != nil {
		l.errs = append(l.errs, "RATE_LIMIT_WINDOW_MS must be a number")
		l.cfg.RateLimitWindowMs = 60000
	} else {
		l.cfg.RateLimitWindowMs = val
	}
}

func (l *loader) getRequired(key string) string {
	val := os.Getenv(key)
	if val == "" {
		l.errs = append(l.errs, fmt.Sprintf("%s is required", key))
	}
	return val
}
