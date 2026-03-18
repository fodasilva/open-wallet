package infra

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/joho/godotenv"
)

type ConfigRoot struct {
	Environment        string
	Delay              string
	Origins            string
	GcpProjectID       string
	Port               string
	DatabaseURL        string
	GoogleClientID     string
	GoogleClientSecret string
	LoginRedirectURI   string
	JWTSecret          string
	RateLimitDBURL     string
	RateLimitMax       string
	RateLimitWindow    string
}

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

var AppConfig *Config

func init() {
	envFiles := []string{".env", ".env.dev", "../../.env", "../../.env.dev"}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			log.Printf("Loaded configuration from %s\n", file)
			break
		}
	}

	ConfigRoot := loadConfig()
	AppConfig = validateConfig(ConfigRoot)

	v := reflect.ValueOf(AppConfig).Elem()
	t := v.Type()

	var builder strings.Builder

	builder.WriteString("-----Environment variables-----\n")

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		valueField := v.Field(i)

		var value any

		if valueField.Kind() == reflect.Ptr {
			if valueField.IsNil() {
				value = nil
			} else {
				value = valueField.Elem().Interface()
			}
		} else {
			value = valueField.Interface()
		}

		builder.WriteString(fmt.Sprintf("%s: %v\n", field.Name, value))
	}

	log.Printf("\n%s", builder.String())
}

func loadConfig() *ConfigRoot {
	return &ConfigRoot{
		Environment:        os.Getenv("ENVIRONMENT"),
		Delay:              os.Getenv("DELAY"),
		Origins:            os.Getenv("ORIGINS"),
		GcpProjectID:       os.Getenv("GCP_PROJECT_ID"),
		Port:               os.Getenv("PORT"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		LoginRedirectURI:   os.Getenv("LOGIN_REDIRECT_URI"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		RateLimitDBURL:     os.Getenv("RATE_LIMIT_DB_URL"),
		RateLimitMax:       os.Getenv("RATE_LIMIT_MAX_REQUESTS"),
		RateLimitWindow:    os.Getenv("RATE_LIMIT_WINDOW_MS"),
	}
}

func validateConfig(ctg *ConfigRoot) *Config {
	errors := make([]string, 0)
	Config := &Config{}

	if ctg.Environment != "" {
		if !slices.Contains([]string{"dev", "prod"}, ctg.Environment) {
			errors = append(errors, "ENVIRONMENT must be 'dev' or 'prod'")
		} else {
			Config.Environment = ctg.Environment
		}
	} else {
		Config.Environment = "dev"
	}

	if ctg.Delay != "" {
		intDelay, err := strconv.Atoi(ctg.Delay)

		if err != nil {
			errors = append(errors, "DELAY must be a number")
		} else {
			Config.Delay = intDelay
		}
	}

	if ctg.Origins != "" {
		splittedOrigins := strings.Split(ctg.Origins, ",")
		errs := make([]string, 0)

		for _, origin := range splittedOrigins {
			if !utils.IsValidURL(origin) {
				errors = append(errs, fmt.Sprintf("ORIGIN %s must be a valid url", origin))
			}
		}

		if len(errs) == 0 {
			Config.Origins = splittedOrigins
		}
	}

	if ctg.GcpProjectID != "" {
		Config.GcpProjectID = &ctg.GcpProjectID
	}

	if ctg.Port != "" {
		intPort, err := strconv.Atoi(ctg.Port)

		if err != nil {
			errors = append(errors, "PORT must be a number")
		} else {
			Config.Port = intPort
		}
	} else {
		Config.Port = 8080
	}

	if ctg.DatabaseURL != "" {
		Config.DatabaseURL = ctg.DatabaseURL
	} else {
		errors = append(errors, "DATABASE_URL is required")
	}

	if ctg.GoogleClientID != "" {
		Config.GoogleClientID = ctg.GoogleClientID
	} else {
		errors = append(errors, "GOOGLE_CLIENT_ID is required")
	}

	if ctg.GoogleClientSecret != "" {
		Config.GoogleSecret = ctg.GoogleClientSecret
	} else {
		errors = append(errors, "GOOGLE_CLIENT_SECRET is required")
	}

	if ctg.LoginRedirectURI != "" {
		Config.LoginRedirectURI = ctg.LoginRedirectURI
	} else {
		errors = append(errors, "LOGIN_REDIRECT_URI is required")
	}

	if ctg.JWTSecret != "" {
		Config.JWTSecret = ctg.JWTSecret
	} else {
		errors = append(errors, "JWT_SECRET is required")
	}

	if ctg.RateLimitDBURL != "" {
		Config.RateLimitDBURL = ctg.RateLimitDBURL
	} else {
		errors = append(errors, "RATE_LIMIT_DB_URL is required")
	}

	if ctg.RateLimitMax != "" {
		intRate, err := strconv.Atoi(ctg.RateLimitMax)
		if err != nil {
			errors = append(errors, "RATE_LIMIT_MAX_REQUESTS must be a number")
		} else {
			Config.RateLimitMaxRequests = intRate
		}
	} else {
		Config.RateLimitMaxRequests = 100
	}

	if ctg.RateLimitWindow != "" {
		intTime, err := strconv.Atoi(ctg.RateLimitWindow)
		if err != nil {
			errors = append(errors, "RATE_LIMIT_WINDOW_MS must be a number")
		} else {
			Config.RateLimitWindowMs = intTime
		}
	} else {
		Config.RateLimitWindowMs = 60000
	}

	if len(errors) > 0 {
		panic("Invalid configuration -> " + strings.Join(errors, ", "))
	}

	return Config
}
