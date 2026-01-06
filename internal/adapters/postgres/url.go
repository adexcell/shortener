package postgres

import (
	"context"

	"github.com/wb-go/wbf/dbpg"
)

type URLPostgres struct {
	db *dbpg.DB
}

func NewURLPostgres(db *dbpg.DB) *URLPostgres {
	return &URLPostgres{db: db}
}

func (p *URLPostgres) Save(ctx context.Context, shortCode, longURL string) error {
	query := `
	INSERT INTO urls (short_code, long_url)
	VALUES ($1, $2)`

	_, err := p.db.ExecContext(ctx, query, shortCode, longURL)
	return err
}

func (p *URLPostgres) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	var longURL string
	query := `
	SELECT long_url FROM urls
	WHERE short_code = $1`
	err := p.db.QueryRowContext(ctx, query, shortCode).Scan(&longURL)
	return longURL, err
}

func (p *URLPostgres) SaveClick(ctx context.Context, shortCode, ip, userAgent string) error {
	query := `
	INSERT INTO analytics (short_code, ip, user_agent)
	VALUES ($1, $2, $3)`
	_, err := p.db.ExecContext(ctx, query, shortCode, ip, userAgent)
	return err
}

func (p *URLPostgres) GetAnalytics(ctx context.Context, shortCode string) (int, error) {
	var count int
	query := `
	SELECT COUNT(*) FROM analytics
	WHERE short_code = $1`
	err := p.db.QueryRowContext(ctx, query, shortCode).Scan(&count)
	return count, err
}
