package storage

import "github.com/blagorodov/go-shortener/internal/app/utils"

type Links map[string]string

var DB Links

// Put Записать короткую ссылку в хранилище
func (l *Links) Put(link string) string {
	key := generateKey()
	DB[key] = link
	return key
}

// Get Получить коротку ссылку из хранилища
func (l *Links) Get(key string) (string, bool) {
	url, ok := DB[key]
	return url, ok
}

// Init Создание хранилища
func Init() {
	DB = make(Links)
}

// Генерация уникального ключа
func generateKey() string {
	var key string
	for {
		key = utils.GenRand(8)
		if _, exists := DB[key]; !exists {
			break
		}
	}
	return key
}
