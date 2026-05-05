package infra

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/felipe1496/open-wallet/internal/util"
)

type Config struct {
	Environment      string
	Delay            int
	Origins          []string
	GcpProjectID     *string
	Port             int
	DatabaseURL      string
	GoogleClientID   string
	GoogleSecret     string
	LoginRedirectURI string
	JWTSecret        string
	RequestTimeout   int
	RateLimits       RateLimits
}

type RateLimits struct {
	XS func() (maxRequests int, windowMs int)
	SM func() (maxRequests int, windowMs int)
	MD func() (maxRequests int, windowMs int)
	LG func() (maxRequests int, windowMs int)
	XL func() (maxRequests int, windowMs int)
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
	l.loadRequestTimeout()
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
			if !util.IsValidURL(part) {
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

func (l *loader) loadRequestTimeout() {
	val := os.Getenv("REQUEST_TIMEOUT_MS")
	if val == "" {
		l.cfg.RequestTimeout = 15000 // default 15s
		return
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		l.errs = append(l.errs, "REQUEST_TIMEOUT_MS must be a number")
		l.cfg.RequestTimeout = 15000
		return
	}
	l.cfg.RequestTimeout = intVal
}

func (l *loader) loadRateLimitConfig() {
	type params struct {
		max    int
		window int
	}

	defaults := map[string]params{
		"XS": {max: 10, window: 60000},
		"SM": {max: 30, window: 60000},
		"MD": {max: 60, window: 60000},
		"LG": {max: 120, window: 60000},
		"XL": {max: 240, window: 60000},
	}

	loadSize := func(size string) (int, int) {
		p := defaults[size]
		maxReqKey := fmt.Sprintf("RATE_LIMIT_%s_MAX_REQUESTS", size)
		windowKey := fmt.Sprintf("RATE_LIMIT_%s_WINDOW_MS", size)

		if val := os.Getenv(maxReqKey); val != "" {
			if intVal, err := strconv.Atoi(val); err == nil {
				p.max = intVal
			} else {
				l.errs = append(l.errs, fmt.Sprintf("%s must be a number", maxReqKey))
			}
		}

		if val := os.Getenv(windowKey); val != "" {
			if intVal, err := strconv.Atoi(val); err == nil {
				p.window = intVal
			} else {
				l.errs = append(l.errs, fmt.Sprintf("%s must be a number", windowKey))
			}
		}
		return p.max, p.window
	}

	xsMax, xsWin := loadSize("XS")
	l.cfg.RateLimits.XS = func() (int, int) { return xsMax, xsWin }

	smMax, smWin := loadSize("SM")
	l.cfg.RateLimits.SM = func() (int, int) { return smMax, smWin }

	mdMax, mdWin := loadSize("MD")
	l.cfg.RateLimits.MD = func() (int, int) { return mdMax, mdWin }

	lgMax, lgWin := loadSize("LG")
	l.cfg.RateLimits.LG = func() (int, int) { return lgMax, lgWin }

	xlMax, xlWin := loadSize("XL")
	l.cfg.RateLimits.XL = func() (int, int) { return xlMax, xlWin }
}

func (l *loader) getRequired(key string) string {
	val := os.Getenv(key)
	if val == "" {
		l.errs = append(l.errs, fmt.Sprintf("%s is required", key))
	}
	return val
}
