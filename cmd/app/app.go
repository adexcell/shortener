package app

import (
	"context"
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
	"github.com/adexcell/shortener/pkg/logger"
	"github.com/adexcell/shortener/pkg/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	cfg     *config.Config
	log     logger.Log
	router  *router.Router
	server  *http.Server
	closers []func() error
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := logger.NewLogger()

	return &App{
		cfg:    cfg,
		log:    log,
		router: router.NewRouter(cfg.Router),
	}, nil
}

func (a *App) Run() {
	a.initDependencies()

	srv := httpserver.New(a.cfg.HTTPServer, a.router)

	// Go routine to start server
	go func() {
		a.log.Info().Str("port", a.cfg.HTTPServer.Port).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error().Err(err).Msg("Server listen failed")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.log.Info().Msg("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	a.shutdown()
}

func (a *App) initDependencies() {
	storage, err := postgres.NewShortenerPostgres(a.cfg.Postgres)
	if err != nil {
		a.log.Fatal().Err(err).Msg("Failed to init Postgres")
	}
	a.addCloser(storage.Close)

	redis := redis.NewShortenerRedis(a.cfg.Redis)
	a.addCloser(redis.Close)

	shortenerUsecase := usecase.NewShortenerUsecase(storage, redis, a.log, a.cfg.Redis.TTL)
	a.addCloser(shortenerUsecase.Close)
	shortenHandler := controller.NewShortenHandler(shortenerUsecase, a.log)

	a.router.Static("/static", "./static")
	a.router.StaticFile("/", "./static/index.html")

	// Swagger
	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	a.log.Info().Msg("register shorten handler")
	shortenHandler.Register(a.router)
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
