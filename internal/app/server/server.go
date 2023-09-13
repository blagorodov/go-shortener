package server

import (
	"github.com/blagorodov/go-shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Start Запуск сервера
func Start() {
	r := chi.NewRouter()
	r.Get("/{key}", handlers.Get)
	r.Post("/", handlers.Post)
	createServer(r)
}

func createServer(r *chi.Mux) {
	if err := http.ListenAndServe(`:8888`, r); err != nil {
		panic(err)
	}
}
