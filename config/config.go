package config

import (
	"github.com/wb-go/wbf/config"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Postgres Postgres `mapstructure:"postgres"`
	Redis    Redis    `mapstructure:"redis"`
}

type App struct {
	Port    string `mapstructure:"port"`
	GinMode string `mapstructure:"gin_mode"`
}

type Postgres struct {
	MasterDSN string   `mapstructure:"master_dsn"`
	SlavesDSN []string `mapstructure:"slaves_dsn"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Load() (*Config, error) {
	cfg := config.New()
	if err := cfg.LoadConfigFiles("config.yaml"); err != nil {
		return nil, err
	}

	var res Config
	if err := cfg.Unmarshal(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
