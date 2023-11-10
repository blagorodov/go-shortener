package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/blagorodov/go-shortener/internal/errs"
	"github.com/blagorodov/go-shortener/internal/models"
	"github.com/blagorodov/go-shortener/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"sync"
	"time"

	"github.com/blagorodov/go-shortener/internal/config"
	def "github.com/blagorodov/go-shortener/internal/repository"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var _ def.Repository = (*Repository)(nil)

type Repository struct {
	m          sync.RWMutex
	connection *sql.DB
	pool       *pgxpool.Pool
}

func NewRepository(ctx context.Context) (*Repository, error) {
	db, err := sql.Open("pgx", config.Options.DBDataSource)
	if err != nil {
		return nil, err
	}

	q := "CREATE TABLE IF NOT EXISTS public.links(key character varying(50) COLLATE pg_catalog.\"default\", link character varying(255) COLLATE pg_catalog.\"default\", CONSTRAINT links_pkey PRIMARY KEY (key))"
	if _, err = db.ExecContext(ctx, q); err != nil {
		return nil, err
	}

	q = "ALTER TABLE public.links ADD COLUMN IF NOT EXISTS user_id character varying(50) COLLATE pg_catalog.\"default\""
	if _, err = db.ExecContext(ctx, q); err != nil {
		return nil, err
	}

	q = "ALTER TABLE public.links ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN DEFAULT FALSE"
	if _, err = db.ExecContext(ctx, q); err != nil {
		return nil, err
	}

	q = "CREATE UNIQUE INDEX IF NOT EXISTS link_unique ON public.links(link)"
	if _, err = db.ExecContext(ctx, q); err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, config.Options.DBDataSource)
	if err != nil {
		return nil, err
	}

	return &Repository{
		connection: db,
		pool:       pool,
	}, nil
}

func (r *Repository) NewKey(ctx context.Context) (string, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var key, dbKey string
	for {
		key = utils.GenRand(8)
		err := r.connection.QueryRowContext(ctx, "SELECT key FROM links WHERE key = $1", key).Scan(&dbKey)
		if errors.Is(err, sql.ErrNoRows) {
			break
		}
		if err != nil {
			return "", err
		}
	}

	return key, nil
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	var url string
	err := r.connection.QueryRowContext(ctx, "SELECT link FROM links WHERE key = $1 AND is_deleted = FALSE", key).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New(errs.ErrKeyNotFound)
	}
	if err != nil {
		return "", err
	}

	return url, nil
}

func (r *Repository) GetKey(ctx context.Context, url string) (string, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	var key string
	err := r.connection.QueryRowContext(ctx, "SELECT key FROM links WHERE link = $1", url).Scan(&key)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New(errs.ErrURLNotFound)
	}
	if err != nil {
		return "", err
	}

	return key, nil
}

func (r *Repository) Put(ctx context.Context, key, url, userID string) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, err := r.pool.Exec(ctx, "INSERT INTO links(key, link, user_id) VALUES($1, $2, $3)", key, url, userID); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Destroy() error {
	if err := r.connection.Close(); err != nil {
		return err
	}
	r.pool.Close()
	return nil
}

func (r *Repository) PingDB(ctx context.Context) error {
	dbCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	return r.connection.PingContext(dbCtx)
}

func (r *Repository) GetURLs(ctx context.Context, userID string) (models.AllResponseList, error) {
	result := make(models.AllResponseList, 0)
	r.m.RLock()
	defer r.m.RUnlock()

	var rows pgx.Rows
	var err error
	if userID != "" {
		rows, err = r.pool.Query(ctx, "SELECT key, link, user_id FROM links WHERE user_id = $1", userID)
	} else {
		rows, err = r.pool.Query(ctx, "SELECT key, link, user_id FROM links")
	}
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var key, link, userID string
		err = rows.Scan(&key, &link, &userID)
		if err != nil {
			return nil, err
		}

		parts := []string{config.Options.BaseURL, key}
		shortURL := strings.Join(parts, `/`)

		var r = models.AllResponse{
			ShortURL:    shortURL,
			OriginalURL: link,
		}
		result = append(result, r)
	}

	return result, nil
}

func (r *Repository) Delete(ctx context.Context, urls []string, userID string) error {
	go deleteURLs(r, ctx, urls, userID)
	return nil
}

func deleteURLs(r *Repository, ctx context.Context, urls []string, userID string) {
	list := make([]string, 0, len(urls))
	// ToDo rewrite with Batch!
	for _, url := range urls {
		// ToDo Need to escape strings!
		list = append(list, fmt.Sprintf("'%s'", url))
	}
	_, err := r.pool.Exec(ctx, "UPDATE links SET is_deleted = TRUE WHERE user_id = $1 AND key IN ($2)", userID, strings.Join(list, ","))
	if err != nil {
		fmt.Println(err)
	}
}
