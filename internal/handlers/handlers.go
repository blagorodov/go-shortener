package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blagorodov/go-shortener/internal/controllers"
	"github.com/blagorodov/go-shortener/internal/cookies"
	"github.com/blagorodov/go-shortener/internal/errs"
	"github.com/blagorodov/go-shortener/internal/models"
	"github.com/blagorodov/go-shortener/internal/service"
	"net/http"
)

// ShortenOne Обработчик POST /api/shorten
func ShortenOne(ctx context.Context, s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ShortenOne Handler")

		cookie, _ := r.Cookie("token")
		if !cookies.Check(r) {
			cookie = cookies.New()
			http.SetCookie(w, cookie)
		}

		userID, _ := cookies.GetIDCookie(cookie)

		url, err := controllers.ShortenOne(ctx, r, s, userID)
		if err != nil && err.Error() != errs.ErrUniqueLinkCode {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
		if err != nil && err.Error() == errs.ErrUniqueLinkCode {
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
		cookie, _ := r.Cookie("token")
		if !cookies.Check(r) {
			cookie = cookies.New()
			http.SetCookie(w, cookie)
		}

		userID, _ := cookies.GetIDCookie(cookie)

		urls, err := controllers.ShortenBatch(ctx, r, s, userID)
		if err != nil && err.Error() != errs.ErrUniqueLinkCode {
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

		if err != nil && err.Error() == errs.ErrUniqueLinkCode {
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
		if !cookies.Check(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}
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

		if !cookies.Check(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		urls, _ := controllers.GetURLs(ctx, r, s)

		result, err := json.Marshal(urls)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
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
