package shortener

import (
	"context"
	"github.com/blagorodov/go-shortener/internal/models"
	"github.com/blagorodov/go-shortener/internal/repository"
)

type Service struct {
	repository repository.Repository
}

func NewService(r repository.Repository) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) NewKey(ctx context.Context) (string, error) {
	return s.repository.NewKey(ctx)
}

func (s *Service) Get(ctx context.Context, key string) (string, error) {
	return s.repository.Get(ctx, key)
}

func (s *Service) GetKey(ctx context.Context, url string) (string, error) {
	return s.repository.GetKey(ctx, url)
}

func (s *Service) Put(ctx context.Context, key, url string, userID int) error {
	return s.repository.Put(ctx, key, url, userID)
}

func (s *Service) PingDB(ctx context.Context) error {
	return s.repository.PingDB(ctx)
}

func (s *Service) GetURLs(ctx context.Context, userID int) (models.AllResponseList, error) {
	return s.repository.GetURLs(ctx, userID)
}
