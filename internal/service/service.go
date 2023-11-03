package service

import (
	"context"
	"github.com/blagorodov/go-shortener/internal/models"
)

type Service interface {
	NewKey(ctx context.Context) (string, error)
	Get(ctx context.Context, key string) (string, error)
	GetKey(ctx context.Context, url string) (string, error)
	Put(ctx context.Context, key, url string) error
	PingDB(ctx context.Context) error
	GetURLs(ctx context.Context, userID int) (models.AllResponseList, error)
}
