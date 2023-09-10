package storage

import "math/rand"

type Links map[string]string

var DB Links

func (l *Links) Put(link string) string {
	key := generateKey()
	DB[key] = link
	return key
}

func (l *Links) Get(key string) (string, bool) {
	url, ok := DB[key]
	return url, ok
}

func Init() {
	DB = make(Links)
}

// Генерация уникального ключа
func generateKey() string {
	var key string
	for {
		key = generateRandomString(8)
		if _, exists := DB[key]; !exists {
			break
		}
	}
	return key
}

// Генерация хэша заданной длины
func generateRandomString(length int) string {
	charset := `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
