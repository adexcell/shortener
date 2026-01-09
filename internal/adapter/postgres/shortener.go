package postgres

import (
	"context"
	"encoding/json"

	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/pkg/postgres"
)

type ShortenerPostgres struct {
	db *postgres.DB
}

func New(cfg postgres.Config) (domain.ShortenerPostgres, error) {
	db, err := postgres.New(cfg)
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
	WITH raw_stats AS (
		-- Шаг 1: Берем все клики по коду один раз
		SELECT 
			clicked_at, 
			user_agent,
			COUNT(*) OVER() as total_count -- считает общее кол-во строк во всем результате
		FROM analytics
		WHERE short_code = $1
	),
	by_date AS (
		-- Шаг 2: Группируем по датам
		SELECT TO_CHAR(clicked_at, 'YYYY-MM-DD') as d, COUNT(*) as c
		FROM raw_stats
		GROUP BY d
		ORDER BY d DESC
		LIMIT 7
	),
	by_browser AS (
		-- Шаг 3: Группируем по браузерам
		SELECT user_agent as b, COUNT(*) as c
		FROM raw_stats
		GROUP BY b
	)
	-- Собираем всё в одну строку
	SELECT 
		COALESCE((SELECT total_count FROM raw_stats LIMIT 1), 0) as total,
		COALESCE((SELECT jsonb_object_agg(d, c) FROM by_date), '{}') as dates,
		COALESCE((SELECT jsonb_object_agg(b, c) FROM by_browser), '{}') as browsers;`

	var dates, browsers []byte
	err := p.db.QueryRowContext(ctx, query, shortCode).Scan(&dto.TotalClicks, &dates, &browsers)
	if err != nil {
		return domain.Stats{}, err
	}

	if err := json.Unmarshal(dates, &dto.ByDate); err != nil {
		return domain.Stats{}, err
	}
	if err := json.Unmarshal(browsers, &dto.ByBrowser); err != nil {
		return domain.Stats{}, err
	}

	res := statsToDomain(dto)

	return res, nil
}

func (p *ShortenerPostgres) Close() error {
	return p.db.Master.Close()
}
