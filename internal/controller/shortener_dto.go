package controller

import (
	"time"

	"github.com/adexcell/shortener/internal/domain"
)

type shortenerControllerDTO struct {
	ID        string    `json:"id"`
	ShortCode string    `json:"short_code"`
	LongURL   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
}

func shortenerToControllerDTO(shortCode, longURL string) (*shortenerControllerDTO, error) {
	s, err := domain.NewShortener(shortCode, longURL)
	if err != nil {
		return &shortenerControllerDTO{}, err
	}

	res := &shortenerControllerDTO{
		ID:        s.ID,
		ShortCode: s.ShortCode,
		LongURL:   s.LongURL,
	}
	return res, nil
}

func shortenerToDomain(dto shortenerControllerDTO) *domain.Shortener {
	return &domain.Shortener{
		ID:        dto.ID,
		ShortCode: dto.ShortCode,
		LongURL:   dto.LongURL,
		CreatedAt: dto.CreatedAt,
	}
}


