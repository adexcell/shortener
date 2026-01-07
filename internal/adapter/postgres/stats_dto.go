package postgres

import (
	"time"

	"github.com/adexcell/shortener/internal/domain"
)

type statsPostgresDTO struct {
	ID          string `db:"id"`
	ShortCode   string `db:"short_code"`
	IP          string `db:"ip"`
	UserAgent   string `db:"user_agent"`
	TotalClicks int
	ByDate      map[string]int
	ByBrowser   map[string]int
	ClickedAt   time.Time `db:"clicked_at"`
}

func statsToPostgresDTO(shortCode, ip, userAgent string) (*statsPostgresDTO, error) {
	s, err := domain.NewStats(shortCode, ip, userAgent)
	if err != nil {
		return &statsPostgresDTO{}, err
	}

	res := &statsPostgresDTO{
		ID:          s.ID,
		ShortCode:   s.ShortCode,
		IP:          s.IP,
		UserAgent:   s.UserAgent,
		TotalClicks: s.TotalClicks,
		ByDate:      s.ByDate,
		ByBrowser:   s.ByBrowser,
		ClickedAt:   s.ClickedAt,
	}
	return res, nil
}

func statsToDomain(dto statsPostgresDTO) domain.Stats {
	return domain.Stats{
		ID:          dto.ID,
		ShortCode:   dto.ShortCode,
		IP:          dto.IP,
		UserAgent:   dto.UserAgent,
		TotalClicks: dto.TotalClicks,
		ByDate:      dto.ByDate,
		ByBrowser:   dto.ByBrowser,
		ClickedAt:   dto.ClickedAt,
	}
}
