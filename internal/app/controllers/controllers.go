package controllers

import (
	"github.com/blagorodov/go-shortener/internal/app/config"
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

// Get Контроллер GET /
func Get(r *http.Request) (string, bool) {
	return storage.Storage.Get(chi.URLParam(r, "id"))
}

// Post Контроллер POST /
func Post(r *http.Request) (string, bool) {
	var url string
	ok := false
	body := readBody(r)

	if len(body) > 0 {
		key := storage.Storage.Put(body)
		parts := []string{config.Options.BaseURL, key}
		url = strings.Join(parts, `/`)
		ok = true
	}
	return url, ok
}

// ReadBody Читаем в строку содержимое Request.Body
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
