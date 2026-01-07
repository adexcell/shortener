package controller

import (
	"time"

	"github.com/adexcell/shortener/internal/domain"
)

type statsControllerDTO struct {
	ID          string         `json:"id"`
	ShortCode   string         `json:"short_code"`
	IP          string         `json:"ip"`
	UserAgent   string         `json:"user_agent"`
	TotalClicks int            `json:"total_clicks"`
	ByDate      map[string]int `json:"by_date"`
	ByBrowser   map[string]int `json:"by_browser"`
	ClickedAt   time.Time      `json:"clicked_at"`
}

func statsToControllerDTO(shortCode, ip, userAgent string) (*statsControllerDTO, error) {
	s, err := domain.NewStats(shortCode, ip, userAgent)
	if err != nil {
		return &statsControllerDTO{}, err
	}

	res := &statsControllerDTO{
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

func statsToResponse(s domain.Stats) statsControllerDTO {
	return statsControllerDTO{
		ID:          s.ID,
		ShortCode:   s.ShortCode,
		IP:          s.IP,
		UserAgent:   s.UserAgent,
		TotalClicks: s.TotalClicks,
		ByDate:      s.ByDate,
		ByBrowser:   s.ByBrowser,
		ClickedAt:   s.ClickedAt,
	}
}
