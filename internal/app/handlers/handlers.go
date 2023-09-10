package handlers

import (
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"io"
	"net/http"
	"strings"
)

// HandleRoot Обработчик /
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		HandlePost(w, r)
	}
	if r.Method == http.MethodGet {
		HandleGet(w, r)
	}
}

// HandlePost Обработчик всех POST /
func HandlePost(w http.ResponseWriter, r *http.Request) {
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

// HandleGet Обработчик всех GET /
func HandleGet(w http.ResponseWriter, r *http.Request) {
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
		url, ok = storage.DB.Get(parts[1])
	}
	return url, ok
}

// Контроллер POST /
func doPost(r *http.Request) (string, bool) {
	var url string
	ok := false
	body := readBody(r)

	if len(body) > 0 {
		key := storage.DB.Put(body)
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
