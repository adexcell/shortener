package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/adexcell/shortener/config"
	_ "github.com/adexcell/shortener/docs" // Swagger docs
	"github.com/adexcell/shortener/internal/adapter/postgres"
	"github.com/adexcell/shortener/internal/adapter/redis"
	"github.com/adexcell/shortener/internal/controller"
	"github.com/adexcell/shortener/internal/usecase"
	"github.com/adexcell/shortener/pkg/httpserver"
	"github.com/adexcell/shortener/pkg/log"
	"github.com/adexcell/shortener/pkg/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	cfg     *config.Config
	log     log.Log
	router  *router.Router
	server  *http.Server
	closers []func() error
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := log.New()

	return &App{
		cfg:    cfg,
		log:    log,
		router: router.New(cfg.Router),
	}, nil
}

func (a *App) Run() error {
	if err := a.initDependencies(); err != nil {
		return err
	}

	srv := httpserver.New(a.router, a.cfg.HTTPServer, a.log)
	a.addCloser(srv.Close)

	// Go routine to start server
	srv.Start()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.log.Info().Msg("Shutting down server...")
	a.shutdown()

	return nil
}

func (a *App) initDependencies() error {
	storage, err := postgres.New(a.cfg.Postgres)
	if err != nil {
		return fmt.Errorf("Failed to init Postgres: %w", err)
	}
	a.addCloser(storage.Close)

	redis := redis.New(a.cfg.Redis)
	a.addCloser(redis.Close)

	shortenerUsecase := usecase.New(storage, redis, a.log, a.cfg.Redis.TTL)
	a.addCloser(shortenerUsecase.Close)
	shortenHandler := controller.NewShortenHandler(shortenerUsecase, a.log)

	a.router.Static("/static", "./static")
	a.router.StaticFile("/", "./static/index.html")

	// Swagger
	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	a.log.Info().Msg("register shorten handler")
	shortenHandler.Register(a.router)

	return nil
}

func (a *App) addCloser(closer func() error) {
	a.closers = append(a.closers, closer)
}

func (a *App) shutdown() {
	for i := len(a.closers) - 1; i >= 0; i-- {
		if err := a.closers[i](); err != nil {
			a.log.Error().Err(err).Msg("failed to close resource")
		}
	}
}
