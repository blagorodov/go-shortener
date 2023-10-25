package file

import (
	"context"
	"github.com/blagorodov/go-shortener/internal/config"
	"strconv"
	"sync"

	def "github.com/blagorodov/go-shortener/internal/repository"
	"github.com/blagorodov/go-shortener/internal/repository/memory"
)

var _ def.Repository = (*Repository)(nil)

type Repository struct {
	memory *memory.Repository
	m      sync.RWMutex
}

func NewRepository(ctx context.Context) (*Repository, error) {
	rp, err := memory.NewRepository(ctx)
	if err != nil {
		return nil, err
	}
	r := &Repository{
		memory: rp,
	}

	urls, err := loadFromFile()
	if err != nil {
		return nil, err
	}

	for _, url := range urls {
		if err := r.memory.Put(ctx, url.ShortURL, url.OriginalURL); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Repository) NewKey(ctx context.Context) (string, error) {
	return r.memory.NewKey(ctx)
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.memory.Get(ctx, key)
}

func (r *Repository) GetKey(ctx context.Context, url string) (string, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.memory.GetKey(ctx, url)
}

func (r *Repository) Put(ctx context.Context, key, url string) error {
	r.m.Lock()
	defer r.m.Unlock()
	if err := r.memory.Put(ctx, key, url); err != nil {
		return err
	}
	return saveToFile(key, url)
}

func (r *Repository) PingDB(_ context.Context) error {
	return nil
}

func (r *Repository) Destroy() error {
	return nil
}

// LoadFromFile Загрузить список ссылок из файла хранилища
func loadFromFile() ([]ShortenURL, error) {
	c, err := newConsumer(config.Options.URLDBPath)
	if err != nil {
		return nil, err
	}
	defer c.close()

	var items []ShortenURL
	for {
		if item, _ := c.readItem(); item != nil {
			uuid, err := strconv.Atoi(item.UUID)
			if err != nil {
				uuid = 1
			}
			lastUUID = uuid
			items = append(items, *item)
		} else {
			break
		}
	}

	return items, nil
}

// SaveToFile Добавить одну ссылку в файл хранилище
func saveToFile(key, link string) error {
	p, err := newProducer(config.Options.URLDBPath)
	if err != nil {
		return err
	}
	defer p.close()

	lastUUID++
	item := &ShortenURL{
		UUID:        strconv.Itoa(lastUUID),
		ShortURL:    key,
		OriginalURL: link,
	}
	return p.writeItem(item)
}
