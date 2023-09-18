package storage

import (
	"math/rand"
)

type linksMap map[string]string

type MemoryStorage struct {
	Storage
	links linksMap
}

// Put Записать короткую ссылку в хранилище
func (l *MemoryStorage) Put(link string) string {
	var key string
	for {
		key = genRand(8)
		if _, exists := l.links[key]; !exists {
			break
		}
	}
	l.links[key] = link
	return key
}

// Get Получить короткую ссылку из хранилища
func (l *MemoryStorage) Get(key string) (string, bool) {
	url, ok := l.links[key]
	return url, ok
}

func NewMemoryStorage() *MemoryStorage {
	r := &MemoryStorage{}
	r.links = make(linksMap)
	return r
}

// GenRand Генерация хэша заданной длины
func genRand(length int) string {
	charset := `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
