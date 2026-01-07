package app

import (
	"github.com/adexcell/shortener/config"
	"github.com/adexcell/shortener/internal/adapter/postgres"
	"github.com/adexcell/shortener/internal/adapter/redis"
	"github.com/adexcell/shortener/internal/controller"
	"github.com/adexcell/shortener/internal/usecase"
	"github.com/adexcell/shortener/pkg/logger"
	"github.com/adexcell/shortener/pkg/router"
)

type App struct {
	log    logger.Log
	cfg    *config.Config
	router *router.Router
}

func NewApp() *App {
	log := logger.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	log.Info().Msg("create router")
	router := router.NewRouter(cfg.Router)

	return &App{
		log:    log,
		cfg:    cfg,
		router: router,
	}
}

func (a *App) Run() {
	a.init()

	a.log.Info().Str("port", a.cfg.HTTPServer.Port).Msg("Starting server")
	if err := a.router.Run(":" + a.cfg.HTTPServer.Port); err != nil {
		a.log.Fatal().Err(err).Msg("Server failed")
	}
}

func (a *App) init() {
	storage, err := postgres.NewShortenerPostgres(a.cfg.Postgres)
	if err != nil {
		a.log.Fatal().Err(err).Msg("Failed to init Postgres")
	}

	rdb := redis.NewShortenerRedis(a.cfg.Redis)

	service := usecase.NewShortenerUsecase(storage, rdb, a.log)
	shortenHandler := controller.NewShortenHandler(service, a.log)

	a.router.Static("/static", "./static")
	a.router.StaticFile("/", "./static/index.html")

	a.log.Info().Msg("register shorten handler")
	shortenHandler.Register(a.router)
}
