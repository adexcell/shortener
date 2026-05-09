package usecase

import (
	"context"
	"fmt"

	"github.com/adexcell/shortener/internal/shortener/domain"
	"github.com/adexcell/shortener/internal/shortener/dto"
)

func (u *ShortenUsecase) GetAnalytics(ctx context.Context, input dto.GetAnalyticsInput) (dto.GetAnalyticsOutput, error) {
	stats := domain.Stats{
		ShortCode: input.ShortCode,
	}

	stats, err := u.postgres.GetDetailedStats(ctx, stats)
	if err != nil {
		return dto.GetAnalyticsOutput{}, fmt.Errorf("postgres.GetDetailedStats: %w", err)
	}

	output := dto.GetAnalyticsOutput{
		TotalClicks: stats.TotalClicks,
		ByDate:      stats.ByDate,
		ByBrowser:   stats.ByBrowser,
	}

	return output, nil
}
