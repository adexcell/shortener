package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/wb-go/wbf/ginext"

	"github.com/adexcell/shortener/config"
	httpserver "github.com/adexcell/shortener/pkg/http/server"
	"github.com/adexcell/shortener/pkg/metrics"
	"github.com/adexcell/shortener/pkg/postgres"
	"github.com/adexcell/shortener/pkg/redis"
)

type Dependencies struct {
	// Adapters
	Postgres *postgres.Pool
	Redis    *redis.Client

	// Controllers
	RouterHTTP *ginext.Engine

	Metrics *metrics.HTTPServer
}

func Run(ctx context.Context, cfg config.Config) (err error) {
	var deps Dependencies

	// Adapters
	deps.Postgres, err = postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	deps.Redis, err = redis.New(ctx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("redis.New: %w", err)
	}

	// Controllers
	deps.RouterHTTP = ginext.New(cfg.Router)

	// Metrics
	deps.Metrics = metrics.NewHTTPServer()

	// Domains
	shorten := ShortenDomain(ctx, deps)

	// Start http server
	httpserver := httpserver.New(deps.RouterHTTP, cfg.HTTP)
	log.Info().Msg("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig // wait signal
	log.Info().Msg("App got signal to stop")

	shorten.Stop()

	// Controllers close
	httpserver.Close()

	// Adapters close
	deps.Redis.Close()
	deps.Postgres.Close()

	log.Info().Msg("App stopped!")

	return nil
}
