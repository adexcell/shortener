package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/adexcell/shortener/internal/shortener/domain"
	"github.com/jackc/pgx/v5"
)

// GetDetailedStats retrieves an analytics about short URL usage.
func (p *Postgres) GetDetailedStats(ctx context.Context, stats domain.Stats) (domain.Stats, error) {

	stats.ByDate = make(map[string]int)
	stats.ByBrowser = make(map[string]int)

	// total clicks
	const sql = `
	WITH raw_stats AS (
		-- Шаг 1: Берем все клики по коду один раз
		SELECT 
			clicked_at, 
			user_agent,
			COUNT(*) OVER() as total_count -- считает общее кол-во строк во всем результате
		FROM clicks
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
	err := p.pgpool.QueryRow(ctx, sql, stats.ShortCode).Scan(&stats.TotalClicks, &dates, &browsers)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Stats{}, domain.ErrShortNotFound
		}
		return domain.Stats{}, err
	}

	if err := json.Unmarshal(dates, &stats.ByDate); err != nil {
		return domain.Stats{}, err
	}
	if err := json.Unmarshal(browsers, &stats.ByBrowser); err != nil {
		return domain.Stats{}, err
	}

	return stats, nil
}
