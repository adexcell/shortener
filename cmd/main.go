package main

import (
	"github.com/adexcell/shortener/config"
	"github.com/adexcell/shortener/internal/adapters/postgres"
	"github.com/adexcell/shortener/internal/controller"
	"github.com/adexcell/shortener/internal/usecase"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.InitConsole()

	cfg, err := config.Load()
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to load config")
	}

	dbOpts := &dbpg.Options{
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	}

	db, err := dbpg.New(cfg.Postgres.MasterDSN, cfg.Postgres.SlavesDSN, dbOpts)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("DB connection failed")
	}
	if err := db.Master.Ping(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("DB Ping failed - check your DSN and SSL mode")
	}
	
	zlog.Logger.Info().Msg("Postgres connected and pinged")

	zlog.Logger.Info().Msg("Postgres connected")

	storage := postgres.NewURLPostgres(db)
	rdb := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	service := usecase.NewShortenerUsecase(storage, rdb)
	shortenHandler := controller.NewShortenHandler(service)

	zlog.Logger.Info().Msg("create router")
	r := ginext.New(cfg.App.GinMode)

	zlog.Logger.Info().Msg("register shorten handler")
	shortenHandler.Register(r)

	zlog.Logger.Info().Str("port", cfg.App.Port).Msg("Starting server")
	if err := r.Run(":" + cfg.App.Port); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Server failed")
	}
}
