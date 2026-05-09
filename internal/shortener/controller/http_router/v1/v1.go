package v1

import (
	"context"

	"github.com/adexcell/shortener/internal/shortener/dto"
)

type Usecase interface {
	CreateShorten(ctx context.Context, input dto.CreateShortenInput) (dto.CreateShortenOutput, error)
	GetOriginalURL(ctx context.Context, input dto.RedirectInput) (string, error)
	GetAnalytics(ctx context.Context, input dto.GetAnalyticsInput) (dto.GetAnalyticsOutput, error)
}

type Handler struct {
	usecase Usecase
}

func New(usecase Usecase) *Handler {
	return &Handler{usecase: usecase}
}
