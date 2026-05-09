package dto

import "time"

// RedirectInput represents the input data required to redirect.
type RedirectInput struct {
	ShortenCode string
	IP          string
	UserAgent   string
	ClickedAt   time.Time
}
