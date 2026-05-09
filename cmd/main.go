package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/adexcell/shortener/config"
	"github.com/adexcell/shortener/internal/app"
	"github.com/adexcell/shortener/pkg/logger"
	"github.com/adexcell/shortener/pkg/otel"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New")
	}

	logger.Init(c.Logger)

	ctx := context.Background()

	if err = otel.Init(ctx, c.OTEL); err != nil {
		log.Error().Err(err).Msg("otel.Init")
	}
	defer otel.Close()

	if err := app.Run(ctx, c); err != nil {
		log.Error().Err(err).Msg("app.Run")
	}
}
