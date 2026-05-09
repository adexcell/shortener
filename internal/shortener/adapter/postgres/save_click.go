package postgres

import (
	"context"

	"github.com/adexcell/shortener/internal/shortener/domain"
)

// SaveClick inserts an analytics datas into database.
func (p *Postgres) SaveClick(ctx context.Context, stats domain.Stats) error {
	const sql = `INSERT INTO clicks (id, short_code, clicked_at, user_agent, ip_address)
			     VALUES ($1, $2, $3, $4, $5)`

	_, err := p.pgpool.Exec(
		ctx,
		sql,
		stats.ID.String(),
		stats.ShortCode,
		stats.ClickedAt,
		stats.UserAgent,
		stats.IP,
	)

	return err
}
