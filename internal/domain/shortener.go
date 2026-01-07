// Package domain содержит основные бизнес модели и интерфейсы.
package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/adexcell/shortener/pkg/utils/uuid"
	"github.com/go-playground/validator/v10"
)

type Shortener struct {
	ID        string
	ShortCode string
	LongURL   string `validate:"required,url"`
	CreatedAt time.Time
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func NewShortener(shortCode, longURL string) (Shortener, error) {
	s := Shortener{
		ID:        uuid.New(),
		ShortCode: shortCode,
		LongURL:   longURL,
	}

	if err := s.Validate(); err != nil {
		return Shortener{}, fmt.Errorf("u.Validate: %w", err)
	}

	return s, nil
}

func (s Shortener) Validate() error {
	err := validate.Struct(s)
	if err != nil {
		return fmt.Errorf("validate.Struct Shortener: %w", err)
	}

	return nil
}

type ShortenerPostgres interface {
	Save(ctx context.Context, shortCode, longURL string) error
	GetLongURL(ctx context.Context, shortCode string) (string, error)
	SaveClick(ctx context.Context, shortCode, ip, userAgent string) error
	GetDetailedStats(ctx context.Context, shortCode string) (Stats, error)
}

type ShortenerRedis interface {
	SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type ShortenerUsecase interface {
	Shorten(ctx context.Context, longURL, alias string) (string, error)
	GetOriginal(ctx context.Context, shortCode, ip, userAgent string) (string, error)
	GetStats(ctx context.Context, shortCode string) (Stats, error)
}
