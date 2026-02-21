package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ConfigRoot struct {
	Enviroment   string
	Delay        string
	Origins      string
	GcpProjectID string
}

type Config struct {
	Enviroment   string
	Delay        int
	Origins      []string
	GcpProjectID *string
}

var AppConfig *Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	ConfigRoot := loadConfig()
	AppConfig = validateConfig(ConfigRoot)
	log.Printf("Enviroment variables loaded: %+v\n", AppConfig)
}

func loadConfig() *ConfigRoot {
	return &ConfigRoot{
		Enviroment:   os.Getenv("ENVIROMENT"),
		Delay:        os.Getenv("DELAY"),
		Origins:      os.Getenv("ORIGINS"),
		GcpProjectID: os.Getenv("GCP_PROJECT_ID"),
	}
}

func validateConfig(ctg *ConfigRoot) *Config {
	errors := make([]string, 0)
	Config := &Config{}

	if ctg.Enviroment != "" {
		if !Contains([]string{"dev", "prod"}, ctg.Enviroment) {
			errors = append(errors, "ENVIROMENT must be 'dev' or 'prod'")
		} else {
			Config.Enviroment = ctg.Enviroment
		}
	} else {
		Config.Enviroment = "dev"
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
			if !isValidURL(origin) {
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

	if len(errors) > 0 {
		panic("Invalid configuration -> " + strings.Join(errors, ", "))
	}

	return Config
}
