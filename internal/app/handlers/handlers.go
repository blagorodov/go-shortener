package handlers

import (
	"github.com/blagorodov/go-shortener/internal/app/controllers"
	"net/http"
)

// Root Обработчик для роута /
func Root(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Post(w, r)
	}
	if r.Method == http.MethodGet {
		Get(w, r)
	}
}

// Post Обработчик всех POST-запросов
func Post(w http.ResponseWriter, r *http.Request) {
	url, ok := controllers.Post(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte(url))
	if err != nil {
		return
	}
}

// Get Обработчик всех GET-запросов
func Get(w http.ResponseWriter, r *http.Request) {
	url, ok := controllers.Get(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set(`Location`, url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
