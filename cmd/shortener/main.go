package main

import (
	"github.com/blagorodov/go-shortener/internal/app/logger"
	"github.com/blagorodov/go-shortener/internal/app/repository"
	"github.com/blagorodov/go-shortener/internal/app/server"
	"github.com/blagorodov/go-shortener/internal/app/storage"
)

func main() {
	logger.Init()
	repository.InitDB()
	defer repository.CloseDB()

	s, err := storage.NewMemoryStorage()
	if err != nil {
		panic(err)
	}
	server.Start(s)
}
