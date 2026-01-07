// Package redis является оберткой над вспомогательным пакетом wbf/redis.
package redis

import "github.com/wb-go/wbf/redis"

type RDB = redis.Client

type Config struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}


func NewRedis(cfg Config) *RDB{
	return redis.New(cfg.Addr, cfg.Password, cfg.DB)
}
