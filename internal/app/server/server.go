package server

import (
	"github.com/blagorodov/go-shortener/internal/app/compress"
	"github.com/blagorodov/go-shortener/internal/app/config"
	"github.com/blagorodov/go-shortener/internal/app/handlers"
	"github.com/blagorodov/go-shortener/internal/app/logger"
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func router(s storage.Storage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use(compress.GzipMiddleware)
	r.Get("/{id}", handlers.Get(s))
	r.Post("/", handlers.Post(s))
	r.Post("/api/shorten", handlers.Post(s))
	return r
}

// Start Запуск сервера
func Start(s storage.Storage) {
	if err := http.ListenAndServe(config.Options.ServerAddress, router(s)); err != nil {
		panic(err)
	}
}
