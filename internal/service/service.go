package service

import (
	"context"
	"time"

	"github.com/adexcell/shortener.git/internal/domain"
	"github.com/adexcell/shortener.git/pkg/keygen"
)

type ShortenerService interface {
	Shorten(ctx context.Context, originalURL string) (string, error)
	Resolve(ctx context.Context, shortCode string) (string, error)
}

type Repository interface {
	SaveLink(ctx context.Context, link domain.Link) error
	GetLink(ctx context.Context, shortCode string) (domain.Link, error)
	DeleteLink(ctx context.Context, shortCode string) error
	SaveStat(ctx context.Context, stat domain.Stat) error
	GetStat(ctx context.Context, linkID string) ([]domain.Stat, error)
}

type Shortener struct {
	repo Repository
}

func NewService(repo Repository) *Shortener {
	return &Shortener{repo: repo}
}

func (s *Shortener) Shorten(ctx context.Context, originalURL string) (string, error) {
	code := keygen.GenerateRandomCode(6)

	link := domain.Link{
		OriginalURL: originalURL,
		ShortCode:   code,
		CreatedAt:   time.Now(),
	}

	err := s.repo.SaveLink(ctx, link)

	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *Shortener) Resolve(ctx context.Context, shortCode string) (string, error) {
	link, err := s.repo.GetLink(ctx, shortCode)
	if err != nil {
		return "", err
	}

	return link.OriginalURL, nil
}
