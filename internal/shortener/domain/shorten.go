package domain

import (
	"fmt"
	"time"
	
	"github.com/google/uuid"

	"github.com/adexcell/shortener/internal/shortener/domain/validation"
)

// Shortener represents a shortener entity in the system.
type Shortener struct {
	ID          uuid.UUID
	ShortCode   string `validate:"lte=20"`
	OriginalURL string `validate:"required,url"`
	CreatedAt   time.Time
}

// NewShortener creates a new Shortener entity with given parameters and validates it.
func NewShortener(shortCode, originalURL string) (Shortener, error) {
	shorten := Shortener{
		ID:          uuid.New(),
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now().UTC(),
	}

	err := validation.Validate(shorten)
	if err != nil {
		return Shortener{}, fmt.Errorf("shorten.Validate: %w", err)
	}

	return shorten, nil
}
