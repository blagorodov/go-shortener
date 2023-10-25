package repository

import "context"

type Repository interface {
	NewKey(ctx context.Context) (string, error)
	Get(ctx context.Context, key string) (string, error)
	GetKey(ctx context.Context, url string) (string, error)
	Put(ctx context.Context, key, url string) error
	PingDB(ctx context.Context) error
	Destroy() error
}
