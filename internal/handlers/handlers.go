package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blagorodov/go-shortener/internal/auth"
	"github.com/blagorodov/go-shortener/internal/controllers"
	"github.com/blagorodov/go-shortener/internal/errs"
	"github.com/blagorodov/go-shortener/internal/logger"
	"github.com/blagorodov/go-shortener/internal/models"
	"github.com/blagorodov/go-shortener/internal/service"
	"net/http"
)

// ShortenOne Обработчик POST /api/shorten
func ShortenOne(ctx context.Context, s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.ContextKey).(string)
		logger.Log(fmt.Sprintf("ShortenOne for user %s", userID))

		url, err := controllers.ShortenOne(ctx, r, s, userID)

		var result []byte
		if r.Header.Get("Content-Type") == "application/json" {
			w.Header().Set("Content-Type", "application/json")

			var errMarshal error
			result, errMarshal = json.Marshal(models.ShortenResponse{Result: url})
			if errMarshal != nil {
				panic(errMarshal)
			}
		} else {
			result = []byte(url)
		}
		if errors.Is(err, errs.ErrUniqueLinkCode) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		_, err = w.Write(result)
		if err != nil {
			return
		}
	}
}

// ShortenBatch Обработчик POST /api/shorten/batch
func ShortenBatch(ctx context.Context, s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.ContextKey).(string)

		logger.Log(fmt.Sprintf("ShortenBatch for user %s", userID))

		urls, err := controllers.ShortenBatch(ctx, r, s, userID)
		if errors.Is(err, errs.ErrUniqueLinkCode) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var result []byte
		if r.Header.Get("Content-Type") == "application/json" {
			w.Header().Set("Content-Type", "application/json")

			var errMarshal error
			result, errMarshal = json.Marshal(urls)
			if errMarshal != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if errors.Is(err, errs.ErrUniqueLinkCode) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		_, err = w.Write(result)
		if err != nil {
			return
		}
	}
}

// Get Обработчик GET /{id}
func Get(ctx context.Context, s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := controllers.Get(ctx, r, s)
		if err != nil {
			if errors.Is(err, errs.ErrKeyNotFound) {
				w.WriteHeader(http.StatusGone)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
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
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// GetUserURLs Список сокращений пользователя
func GetUserURLs(ctx context.Context, s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value(auth.ContextKey).(string)
		logger.Log(fmt.Sprintf("GetUserURLs for user %s", userID))

		urls, _ := controllers.GetURLs(ctx, s, userID)

		result, err := json.Marshal(urls)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}

		_, err = w.Write(result)
		if err != nil {
			return
		}
	}
}

func Delete(ctx context.Context, s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.ContextKey).(string)
		logger.Log(fmt.Sprintf("Delete for user %s", userID))

		err := controllers.Delete(ctx, r, s, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
