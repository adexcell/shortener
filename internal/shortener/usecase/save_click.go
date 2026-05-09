package usecase

import (
	"context"

	"github.com/adexcell/shortener/internal/shortener/domain"
	"github.com/adexcell/shortener/internal/shortener/dto"
	"github.com/rs/zerolog/log"
)

func (u *ShortenUsecase) SaveClick(ctx context.Context, input dto.RedirectInput) error {
	stats, err := domain.NewStats(input.ShortenCode, input.IP, input.UserAgent, input.ClickedAt)
	if err != nil {
		log.Error().Err(err).
			Str("shorten", input.ShortenCode).
			Str("ip", input.IP).
			Str("user-agent", input.UserAgent).
			Time("clicked at", input.ClickedAt).
			Msg("wrong stats format")
		return err

	}

	err = u.postgres.SaveClick(ctx, stats)
	if err != nil {
		log.Error().Err(err).Msg("failed to save click")
		return err
	}
	
	log.Debug().
		Str("shorten", input.ShortenCode).
		Str("ip", input.IP).
		Str("user-agent", input.UserAgent).
		Time("clicked at", input.ClickedAt).
		Msg("saved stats")

	return nil
}
