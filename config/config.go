package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	httpserver "github.com/adexcell/shortener/pkg/http/server"
	"github.com/adexcell/shortener/pkg/logger"
	"github.com/adexcell/shortener/pkg/otel"
	"github.com/adexcell/shortener/pkg/postgres"
	"github.com/adexcell/shortener/pkg/redis"
)

// App holds basic application metadata.
type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}

// Config represents the complete application configuration.
type Config struct {
	App      App
	HTTP     httpserver.Config
	Logger   logger.Config
	OTEL     otel.Config
	Postgres postgres.Config
	Redis    redis.Config
	Router string `envconfig:"GIN_MODE"`
}

// New loads the configuration from environment variables and .env file.
func New() (Config, error) {
	var config Config

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: .env not loaded: %v", err)
	}

	err = envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
