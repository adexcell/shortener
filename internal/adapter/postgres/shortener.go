package postgres

import (
	"context"

	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/pkg/postgres"
)

type ShortenerPostgres struct {
	db *postgres.DB
}

func NewShortenerPostgres(cfg postgres.Config) (domain.ShortenerPostgres, error) {
	db, err := postgres.NewPostgres(cfg)
	return &ShortenerPostgres{db: db}, err
}

func (p *ShortenerPostgres) Save(ctx context.Context, shortCode, longURL string) error {
	dto, err := shortenerToPostgresDTO(shortCode, longURL)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO urls (id, short_code, long_url)
	VALUES ($1, $2, $3)`

	_, err = p.db.ExecContext(ctx, query, dto.ID, dto.ShortCode, dto.LongURL)
	return err
}

func (p *ShortenerPostgres) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	var longURL string
	query := `
	SELECT long_url FROM urls
	WHERE short_code = $1`
	err := p.db.QueryRowContext(ctx, query, shortCode).Scan(&longURL)
	return longURL, err
}

func (p *ShortenerPostgres) SaveClick(ctx context.Context, shortCode, ip, userAgent string) error {
	dto, err := statsToPostgresDTO(shortCode, ip, userAgent)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO analytics (id, short_code, ip, user_agent)
	VALUES ($1, $2, $3, $4)`

	_, err = p.db.ExecContext(
		ctx,
		query,
		dto.ID,
		dto.ShortCode,
		dto.IP,
		dto.UserAgent,
	)
	return err
}

func (p *ShortenerPostgres) GetDetailedStats(ctx context.Context, shortCode string) (domain.Stats, error) {
	var dto statsPostgresDTO
	dto.ByDate = make(map[string]int)
	dto.ByBrowser = make(map[string]int)

	// total clicks
	query := `
	SELECT COUNT(*) FROM analytics
	WHERE short_code = $1`

	err := p.db.QueryRowContext(ctx, query, shortCode).Scan(&dto.TotalClicks)
	if err != nil {
		return domain.Stats{}, err
	}

	// clicks by date
	query = `
	SELECT TO_CHAR(clicked_at, 'YYYY-MM-DD') as date, COUNT(*)
	FROM analytics
	WHERE short_code = $1
	GROUP BY date
	ORDER BY date DESC
	LIMIT 7`

	rows, err := p.db.QueryContext(ctx, query, shortCode)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date string
			var count int
			if err := rows.Scan(&date, &count); err == nil {
				dto.ByDate[date] = count
			}
		}
	}

	// clicks by browser
	query = `
	SELECT user_agent, COUNT(*)
	FROM analytics
	WHERE short_code = $1
	GROUP BY user_agent`

	rows, _ = p.db.QueryContext(ctx, query, shortCode)
	for rows.Next() {
		var userAgent string
		var count int
		rows.Scan(&userAgent, &count)
		dto.ByBrowser[userAgent] = count
	}

	res := statsToDomain(dto)

	return res, nil
}

func (p *ShortenerPostgres) Close() error {
	return p.db.Master.Close()
}
