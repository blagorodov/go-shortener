package handlers

import (
	"context"
	"encoding/json"
	"github.com/blagorodov/go-shortener/internal/controllers"
	"github.com/blagorodov/go-shortener/internal/models"
	"github.com/blagorodov/go-shortener/internal/service"
	"net/http"
)

// Post Обработчик всех POST-запросов
func Post(ctx context.Context, s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := controllers.Post(ctx, r, s)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var result []byte
		if r.Header.Get("Content-Type") == "application/json" {
			w.Header().Set("Content-Type", "application/json")

			result, err = json.Marshal(models.ShortenResponse{Result: url})
			if err != nil {
				panic(err)
			}
		} else {
			result = []byte(url)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(result)
		if err != nil {
			return
		}
	}
}

// Get Обработчик всех GET-запросов
func Get(ctx context.Context, s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := controllers.Get(ctx, r, s)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set(`Location`, url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// PingDB Проверка подключения к БД
func PingDB(ctx context.Context, s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.PingDB(ctx)
		if err != nil {
			w.WriteHeader(500)
		} else {
			w.WriteHeader(200)
		}
	}
}
