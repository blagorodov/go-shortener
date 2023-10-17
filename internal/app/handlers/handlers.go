package handlers

import (
	"encoding/json"
	"github.com/blagorodov/go-shortener/internal/app/controllers"
	"github.com/blagorodov/go-shortener/internal/app/models"
	"github.com/blagorodov/go-shortener/internal/app/repository"
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"net/http"
)

// Post Обработчик всех POST-запросов
func Post(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := controllers.Post(r, s)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var result []byte
		var err error
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
func Get(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := controllers.Get(r, s)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set(`Location`, url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func PingDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repository.PingDB()
		if err != nil {
			w.WriteHeader(500)
		} else {
			w.WriteHeader(200)
		}
	}
}
