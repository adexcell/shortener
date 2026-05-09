package app

import (
	"context"

	"github.com/adexcell/shortener/internal/shortener/adapter/postgres"
	"github.com/adexcell/shortener/internal/shortener/adapter/redis"
	httprouter "github.com/adexcell/shortener/internal/shortener/controller/http_router"
	"github.com/adexcell/shortener/internal/shortener/usecase"
	"github.com/adexcell/shortener/internal/shortener/worker"
)

type Shorten struct {
	postgres         *postgres.Postgres
	redis            *redis.Redis
	asyncRedisWriter *worker.AsyncRedisWriter
	usecase          *usecase.ShortenUsecase
}

func ShortenDomain(ctx context.Context, deps Dependencies) Shorten{
	var shorten Shorten
	
	shorten.postgres = postgres.New(deps.Postgres)
	shorten.redis = redis.New(deps.Redis)

	shorten.asyncRedisWriter = worker.NewAsyncRedisWriter(shorten.redis)
	shorten.asyncRedisWriter.Start(ctx)

	shorten.usecase = usecase.New(
		ctx,
		shorten.postgres,
		shorten.redis,
		shorten.asyncRedisWriter,
	)

	httprouter.ShortenRouter(deps.RouterHTTP, shorten.usecase, deps.Metrics)
	return shorten
}

func (s Shorten) Stop() {
	s.usecase.Stop()
	s.asyncRedisWriter.Stop()
}
