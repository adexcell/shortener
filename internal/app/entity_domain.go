package app

import (
	"github.com/adexcell/shortener/internal/entity/adapter/postgres"
	"github.com/adexcell/shortener/internal/entity/adapter/rabbit"
	"github.com/adexcell/shortener/internal/entity/adapter/redis"
	httprouter "github.com/adexcell/shortener/internal/entity/controller/http_router"
	"github.com/adexcell/shortener/internal/entity/usecase"
)

func EntityDomain(d Dependencies) {
	entityUseCase := usecase.New(
		postgres.New(d.Postgres),
		redis.New(d.Redis),
		rabbit.New(d.RabbitMQ),
	)

	httprouter.EntityRouter(d.RouterHTTP, entityUseCase, d.Metrics)
}
