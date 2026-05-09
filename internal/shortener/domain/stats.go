package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/adexcell/shortener/internal/shortener/domain/validation"
)

// Stats represents a statistics entity in the system.
type Stats struct {
	ID          uuid.UUID
	ShortCode   string
	IP          string `validate:"required,ip"`
	UserAgent   string
	TotalClicks int
	ByDate      map[string]int
	ByBrowser   map[string]int
	ClickedAt   time.Time
}

// NewStats creates a new Stats entity with the given parameters and validates it.
func NewStats(shortCode, ip, userAgent string, clickedAt time.Time) (Stats, error) {
	stats := Stats{
		ID:        uuid.New(),
		ShortCode: shortCode,
		IP:        ip,
		UserAgent: userAgent,
		ClickedAt: clickedAt,
	}

	err := validation.Validate(stats)
	if err != nil {
		return Stats{}, fmt.Errorf("shorten.Validate: %w", err)
	}

	return stats, nil
}
