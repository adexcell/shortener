package domain

import "time"

type Link struct {
	ID          string    `db:"link_id"`
	OriginalURL string    `db:"original_url"`
	ShortCode   string    `db:"short_code"`
	CreatedAt   time.Time `db:"created_at"`
}

type Stat struct {
	ID        string    `db:"stat_id"`
	LinkID    string    `db:"link_id"`
	Timestamp time.Time `db:"created_at"`
	IP        string    `db:"ip"`
	UserAgent string    `db:"user_agent"`
}
