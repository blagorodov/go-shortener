package main

import (
	"github.com/blagorodov/go-shortener/internal/app/logger"
	"github.com/blagorodov/go-shortener/internal/app/server"
	"github.com/blagorodov/go-shortener/internal/app/storage"
)

func main() {
	logger.Init()
	s, err := storage.NewMemoryStorage()
	if err != nil {
		panic(err)
	}
	server.Start(s)
}
