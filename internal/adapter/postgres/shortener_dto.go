package postgres

import (
	"time"

	"github.com/adexcell/shortener/internal/domain"
)

type shortenerPostgresDTO struct {
	ID        string    `db:"id"`
	ShortCode string    `db:"short_code"`
	LongURL   string    `db:"long_url"`
	CreatedAt time.Time `db:"created_at"`
}

func shortenerToPostgresDTO(shortCode, longURL string) (*shortenerPostgresDTO, error) {
	s, err := domain.NewShortener(shortCode, longURL)
	if err != nil {
		return &shortenerPostgresDTO{}, err
	}

	res := &shortenerPostgresDTO{
		ID:        s.ID,
		ShortCode: s.ShortCode,
		LongURL:   s.LongURL,
	}
	return res, nil
}

func shortenerToDomain(dto shortenerPostgresDTO) *domain.Shortener {
	return &domain.Shortener{
		ID:        dto.ID,
		ShortCode: dto.ShortCode,
		LongURL:   dto.LongURL,
		CreatedAt: dto.CreatedAt,
	}
}


