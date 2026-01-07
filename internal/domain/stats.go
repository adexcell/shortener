package domain

import (
	"fmt"
	"time"

	"github.com/adexcell/shortener/pkg/utils/uuid"
)

type Stats struct {
	ID          string
	ShortCode   string
	IP          string `validate:"required,ip"`
	UserAgent   string
	TotalClicks int
	ByDate      map[string]int
	ByBrowser   map[string]int
	ClickedAt   time.Time
}

func NewStats(shortCode, ip, userAgent string) (Stats, error) {
	s := Stats{
		ID:        uuid.New(),
		ShortCode: shortCode,
		IP:        ip,
		UserAgent: userAgent,
	}

	if err := s.Validate(); err != nil {
		return Stats{}, fmt.Errorf("u.Validate: %w", err)
	}

	return s, nil
}

func (s Stats) Validate() error {
	err := validate.Struct(s)
	if err != nil {
		return fmt.Errorf("validate.Struct Stats: %w", err)
	}

	return nil
}
