package service

import "context"

type Service interface {
	NewKey(ctx context.Context) (string, error)
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key, url string) error
	PingDB(ctx context.Context) error
}
