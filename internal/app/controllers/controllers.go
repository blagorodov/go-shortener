package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/blagorodov/go-shortener/internal/app/config"
	"github.com/blagorodov/go-shortener/internal/app/models"
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

// Get Контроллер GET /
func Get(r *http.Request, s storage.Storage) (string, bool) {
	return s.Get(chi.URLParam(r, "id"))
}

// Post Контроллер POST /
func Post(r *http.Request, s storage.Storage) (string, bool) {
	var url string
	url, ok := getUrl(r)

	if len(url) > 0 {
		key := s.Put(url)
		parts := []string{config.Options.BaseURL, key}
		url = strings.Join(parts, `/`)
		ok = true
	} else {
		ok = false
	}
	return url, ok
}

func getUrl(r *http.Request) (string, bool) {
	if r.Header.Get("Content-Type") == "application/json" {
		var request models.ShortenRequest
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			return "", false
		}
		if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
			return "", false
		}

		return request.Url, true
	} else {
		return readBody(r), true
	}
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
