package server

import (
	"context"
	"github.com/blagorodov/go-shortener/internal/compress"
	"github.com/blagorodov/go-shortener/internal/config"
	"github.com/blagorodov/go-shortener/internal/handlers"
	"github.com/blagorodov/go-shortener/internal/logger"
	"github.com/blagorodov/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	router  *chi.Mux
	service service.Service
}

func NewServer(ctx context.Context, service service.Service) *Server {
	s := Server{
		service: service,
		router:  chi.NewRouter(),
	}
	s.router.Use(logger.WithLogging)
	s.router.Use(compress.GzipMiddleware)
	s.router.Get("/{id}", handlers.Get(ctx, s.service))
	s.router.Get("/ping", handlers.PingDB(ctx, s.service))
	s.router.Post("/", handlers.ShortenOne(ctx, s.service))
	s.router.Post("/api/shorten", handlers.ShortenOne(ctx, s.service))
	s.router.Post("/api/shorten/batch", handlers.ShortenBatch(ctx, s.service))
	s.router.Post("/api/user/login", handlers.Login())
	s.router.Post("/api/user/urls", handlers.GetUserURLs(ctx, s.service))
	return &s
}

// Start Запуск сервера
func (s *Server) Start() {
	if err := http.ListenAndServe(config.Options.ServerAddress, s.router); err != nil {
		panic(err)
	}
}
