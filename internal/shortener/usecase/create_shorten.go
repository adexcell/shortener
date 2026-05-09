package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/adexcell/shortener/internal/shortener/domain"
	"github.com/adexcell/shortener/internal/shortener/dto"
)

func (u *ShortenUsecase) CreateShorten(ctx context.Context, input dto.CreateShortenInput) (dto.CreateShortenOutput, error) {
	shortenCode := input.ShortenCode
	if shortenCode == "" {
		b := make([]byte, 20)
		_, err := rand.Read(b)
		if err != nil {
			return dto.CreateShortenOutput{}, fmt.Errorf("rand.Read: %w", err)
		}
		shortenCode = base64.URLEncoding.EncodeToString(b)[:6]
	}

	shorten, err := domain.NewShortener(shortenCode, input.OriginalURL)
	if err != nil {
		return dto.CreateShortenOutput{}, fmt.Errorf("domain.NewShortener: %w", err)
	}

	err = u.postgres.SaveURL(ctx, shorten)
	if err != nil {
		return dto.CreateShortenOutput{}, fmt.Errorf("postgres.SaveURL: %w", err)
	}

	// save to cache
	task := dto.ShortenTask{
		Shorten:     shorten.ShortCode,
		OriginalURL: input.OriginalURL,
	}
	u.asyncRedisWriter.Send(ctx, task)

	return dto.CreateShortenOutput{ShortenCode: shorten.ShortCode}, nil
}
