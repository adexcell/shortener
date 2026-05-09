package postgres

import (
	"context"
	"errors"

	"github.com/adexcell/shortener/internal/shortener/domain"
	"github.com/jackc/pgx/v5"
)

// GetOriginalURL retrieves an original URL from the database by its shortcode.
func (p *Postgres) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	var originalURL string

	const sql = `SELECT original_url FROM urls
				 WHERE short_code = $1`

	err := p.pgpool.QueryRow(ctx, sql, shortCode).Scan(&originalURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrShortNotFound
	}

	return originalURL, err
}
