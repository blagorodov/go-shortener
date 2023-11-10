package repository

import (
	"context"
	"github.com/blagorodov/go-shortener/internal/models"
)

type Repository interface {
	NewKey(ctx context.Context) (string, error)
	Get(ctx context.Context, key string) (string, error)
	GetKey(ctx context.Context, url string) (string, error)
	Put(ctx context.Context, key, url, userID string) error
	PingDB(ctx context.Context) error
	Destroy() error
	GetURLs(ctx context.Context, userID string) (models.AllResponseList, error)
	Delete(ctx context.Context, urls []string, userID string) error
}
