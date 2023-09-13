package server

import (
	"github.com/blagorodov/go-shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id}", handlers.Get)
	r.Post("/", handlers.Post)
	return r
}

// Start Запуск сервера
func Start() {
	if err := http.ListenAndServe(`:8080`, Router()); err != nil {
		panic(err)
	}
}
