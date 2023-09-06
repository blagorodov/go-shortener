package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type Links map[string]string

var links Links

func main() {
	links = make(Links)
	http.HandleFunc(`/`, handleRoot)
	if err := http.ListenAndServe(`:8889`, nil); err != nil {
		panic(err)
	}
}

// Обработчик /
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handlePost(w, r)
	}
	if r.Method == http.MethodGet {
		handleGet(w, r)
	}
}

// Обработчик всех POST /
func handlePost(w http.ResponseWriter, r *http.Request) {
	url, ok := doPost(r)
	if !ok {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(201)
	_, err := w.Write([]byte(url))
	if err != nil {
		return
	}
}

// Обработчик всех GET /
func handleGet(w http.ResponseWriter, r *http.Request) {
	url, ok := doGet(r)
	if !ok {
		w.WriteHeader(400)
		return
	}
	w.Header().Set(`Location`, url)
	w.WriteHeader(307)
}

// Контроллер GET /
func doGet(r *http.Request) (string, bool) {
	var url string
	ok := false
	parts := strings.Split(r.URL.String(), `/`)
	if len(parts) == 2 {
		url, ok = links[parts[1]]
	}
	return url, ok
}

// Контроллер POST /
func doPost(r *http.Request) (string, bool) {
	var url string
	ok := false
	body := readBody(r)

	if len(body) > 0 {
		key := generateKey()
		links[key] = body
		parts := []string{`http:/`, r.Host, key}
		url = strings.Join(parts, `/`)
		ok = true
	}

	return url, ok
}

// Читаем в строку содержимое Request.Body
func readBody(r *http.Request) string {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)
	body, _ := io.ReadAll(r.Body)
	return string(body)
}

// Генерация уникального ключа
func generateKey() string {
	var key string
	for {
		key = generateRandomString(8)
		if _, exists := links[key]; !exists {
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
