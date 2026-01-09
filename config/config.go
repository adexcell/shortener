package config

import (
	"github.com/adexcell/shortener/pkg/httpserver"
	"github.com/adexcell/shortener/pkg/postgres"
	"github.com/adexcell/shortener/pkg/redis"
	"github.com/adexcell/shortener/pkg/router"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	App        App
	HTTPServer httpserver.Config
	Router     router.Config
	Postgres   postgres.Config
	Redis      redis.Config
}

type App struct {
	Name    string `mapstructure:"app_name"`
	Version string `mapstructure:"app_version"`
}

func Load() (*Config, error) {
	cfg := config.New()

	cfg.EnableEnv("")

	_ = cfg.LoadEnvFiles(".env")

	if err := cfg.LoadConfigFiles("config/config.yaml"); err != nil {
		return nil, err
	}

	var res Config
	if err := cfg.Unmarshal(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
