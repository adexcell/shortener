package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/adexcell/shortener.git/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrLinkNotFound  = errors.New("link not found")
	ErrAlreadyExists = errors.New("short code already exists")
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, connString string) (*Storage, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) SaveLink(ctx context.Context, link domain.Link) error {
	query := `INSERT INTO link (original_url, short_code, created_at) VALUES ($1, $2, $3)`

	_, err := s.pool.Exec(ctx, query, link.OriginalURL, link.ShortCode, link.CreatedAt)
	if err != nil {
		// Check pgErr *pgconn.PgError
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (s *Storage) GetLink(ctx context.Context, shortCode string) (domain.Link, error) {
	query := `SELECT link_id, original_url, short_code, created_at FROM link WHERE short_code = $1`

	row, _ := s.pool.Query(ctx, query, shortCode)
	link, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[domain.Link])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Link{}, ErrLinkNotFound
		}
		return domain.Link{}, err
	}

	return link, nil
}

func (s *Storage) DeleteLink(ctx context.Context, shortCode string) error {
	query := `DELETE FROM link WHERE short_code = $1`

	tag, err := s.pool.Exec(ctx, query, shortCode)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrLinkNotFound
	}
	return nil
}

func (s *Storage) SaveStat(ctx context.Context, stat domain.Stat) error {
	query := `INSERT INTO stat (link_id, created_at, ip, user_agent) VALUES ($1, $2, $3, $4)`

	_, err := s.pool.Exec(ctx, query, stat.LinkID, stat.Timestamp, stat.IP, stat.UserAgent)
	return err
}

func (s *Storage) GetStat(ctx context.Context, linkID string) ([]domain.Stat, error) {
	query := `SELECT link_id, timestamp, ip, user_agent FROM stat WHERE stat_id = $1`

	rows, _ := s.pool.Query(ctx, query, linkID)

	stats, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Stat])
	if err != nil {
		return nil, err
	}

	return stats, nil
}
