package db

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/blagorodov/go-shortener/internal/config"
	def "github.com/blagorodov/go-shortener/internal/repository"
)

var _ def.Repository = (*Repository)(nil)

type Repository struct {
	config     config.Config
	m          sync.RWMutex
	connection *sql.DB
}

func NewRepository(config config.Config) *Repository {
	db, err := sql.Open("pgx", config.DBDataSource)
	if err != nil {
		panic(err)
	}
	return &Repository{
		config:     config,
		connection: db,
	}
}

func (r *Repository) NewKey(ctx context.Context) (string, error) {
	return "", nil
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	return "", nil
}

func (r *Repository) Put(ctx context.Context, key, url string) error {
	r.m.Lock()
	defer r.m.Unlock()
	return nil
}

func (r *Repository) Destroy() error {
	if err := r.connection.Close(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) PingDB(ctx context.Context) error {
	dbCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	return r.connection.PingContext(dbCtx)
}
