package controllers

import (
	"github.com/blagorodov/go-shortener/internal/app/config"
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"github.com/blagorodov/go-shortener/internal/app/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

// Get Контроллер GET /
func Get(r *http.Request) (string, bool) {
	return storage.DB.Get(chi.URLParam(r, "id"))
}

// Post Контроллер POST /
func Post(r *http.Request) (string, bool) {
	var url string
	ok := false
	body := utils.ReadBody(r)

	if len(body) > 0 {
		key := storage.DB.Put(body)
		parts := []string{config.Options.ResultHost, key}
		url = strings.Join(parts, `/`)
		ok = true
	}
	return url, ok
}
