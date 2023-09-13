package controllers

import (
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"github.com/blagorodov/go-shortener/internal/app/utils"
	"net/http"
	"strings"
)

// Get Контроллер GET /
func Get(r *http.Request) (string, bool) {
	var url string
	ok := false
	parts := strings.Split(r.URL.String(), `/`)
	if len(parts) == 2 {
		url, ok = storage.DB.Get(parts[1])
	}
	return url, ok
}

// Post Контроллер POST /
func Post(r *http.Request) (string, bool) {
	var url string
	ok := false
	body := utils.ReadBody(r)

	if len(body) > 0 {
		key := storage.DB.Put(body)
		parts := []string{`http:/`, r.Host, key}
		url = strings.Join(parts, `/`)
		ok = true
	}

	return url, ok
}
