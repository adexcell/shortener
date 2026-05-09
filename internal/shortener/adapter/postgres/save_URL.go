package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/adexcell/shortener/internal/shortener/domain"
)

// SaveURL inserts a new shorten URL into database.
func (p *Postgres) SaveURL(ctx context.Context, shorten domain.Shortener) error {
	const sql = `
	INSERT INTO urls (id, short_code, original_url, created_at)
	VALUES ($1, $2, $3, $4)`

	_, err := p.pgpool.Exec(
		ctx,
		sql,
		shorten.ID.String(),
		shorten.ShortCode,
		shorten.OriginalURL,
		shorten.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Код 23505 означает unique_violation
				return domain.ErrShortCodeAlreadyExists
			}
		}
		return err
	}
	return nil
}
