package memory

import (
	"context"
	"errors"
	"sync"

	def "github.com/blagorodov/go-shortener/internal/repository"
	"github.com/blagorodov/go-shortener/internal/utils"
)

var _ def.Repository = (*Repository)(nil)

type linksMap map[string]string

type Repository struct {
	data linksMap
	m    sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		data: make(linksMap),
	}
}

func (r *Repository) NewKey(_ context.Context) (string, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var key string
	for {
		key = utils.GenRand(8)
		if _, exists := r.data[key]; !exists {
			break
		}
	}
	return key, nil
}

func (r *Repository) Get(_ context.Context, key string) (string, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	url, ok := r.data[key]
	if ok {
		return url, nil
	}
	return "", errors.New("ключ не найден")
}

func (r *Repository) Put(_ context.Context, key, url string) error {
	r.data[key] = url
	return nil
}

func (r *Repository) PingDB(ctx context.Context) error {
	return nil
}

func (r *Repository) Destroy() error {
	return nil
}