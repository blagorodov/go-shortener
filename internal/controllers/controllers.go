package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/blagorodov/go-shortener/internal/config"
	"github.com/blagorodov/go-shortener/internal/models"
	"github.com/blagorodov/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

// Get Контроллер GET /
func Get(ctx context.Context, r *http.Request, s service.Service) (string, error) {
	return s.Get(ctx, chi.URLParam(r, "id"))
}

// Post Контроллер POST /
func Post(ctx context.Context, r *http.Request, s service.Service) (string, error) {
	var url string
	url, err := getURL(r)
	if err != nil {
		return "", err
	}

	if len(url) == 0 {
		return "", errors.New("пустой url")
	}

	key, err := s.NewKey(ctx)
	if err != nil {
		return "", err
	}

	err = s.Put(ctx, key, url)
	if err != nil {
		return "", err
	}

	parts := []string{config.Options.BaseURL, key}
	url = strings.Join(parts, `/`)

	return url, nil
}

func getURL(r *http.Request) (string, error) {
	if r.Header.Get("Content-Type") == "application/json" {
		var request models.ShortenRequest
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			return "", err
		}
		if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
			return "", err
		}

		return request.URL, nil
	} else {
		return readBody(r), nil
	}
}

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
