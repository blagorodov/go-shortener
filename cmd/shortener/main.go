package main

import (
	"github.com/blagorodov/go-shortener/internal/app/server"
	"github.com/blagorodov/go-shortener/internal/app/storage"
)

func main() {
	storage.Init()
	server.Start()
}
