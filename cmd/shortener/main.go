package main

import (
	"github.com/blagorodov/go-shortener/internal/app/config"
	"github.com/blagorodov/go-shortener/internal/app/server"
)

func main() {
	config.ParseFlags()
	server.Start()
}
