package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adexcell/shortener/config"
	_ "github.com/adexcell/shortener/docs" // Swagger docs
	"github.com/adexcell/shortener/internal/adapter/postgres"
	"github.com/adexcell/shortener/internal/adapter/redis"
	"github.com/adexcell/shortener/internal/controller"
	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/internal/usecase"
	"github.com/adexcell/shortener/pkg/logger"
	"github.com/adexcell/shortener/pkg/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	log      logger.Log
	cfg      *config.Config
	router   *router.Router
	postgres domain.ShortenerPostgres
	redis    domain.ShortenerRedis
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

	srv := &http.Server{
		Addr:    ":" + a.cfg.HTTPServer.Port,
		Handler: a.router,
	}

	// Go routine to start server
	go func() {
		a.log.Info().Str("port", a.cfg.HTTPServer.Port).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Fatal().Err(err).Msg("Server listen failed")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.log.Info().Msg("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	// Close database connections
	a.log.Info().Msg("Closing database connections...")
	if err := a.postgres.Close(); err != nil {
		a.log.Error().Err(err).Msg("Failed to close Postgres connection")
	}
	if err := a.redis.Close(); err != nil {
		a.log.Error().Err(err).Msg("Failed to close Redis connection")
	}

	a.log.Info().Msg("Server exiting")
}

func (a *App) init() {
	storage, err := postgres.NewShortenerPostgres(a.cfg.Postgres)
	if err != nil {
		a.log.Fatal().Err(err).Msg("Failed to init Postgres")
	}
	a.postgres = storage

	rdb := redis.NewShortenerRedis(a.cfg.Redis)
	a.redis = rdb

	service := usecase.NewShortenerUsecase(a.postgres, a.redis, a.log)
	shortenHandler := controller.NewShortenHandler(service, a.log)

	a.router.Static("/static", "./static")
	a.router.StaticFile("/", "./static/index.html")

	// Swagger
	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	a.log.Info().Msg("register shorten handler")
	shortenHandler.Register(a.router)
}
