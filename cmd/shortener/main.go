package main

import (
	"github.com/blagorodov/go-shortener/internal/app/handlers"
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"net/http"
)

func main() {
	storage.Init()
	http.HandleFunc(`/`, handlers.HandleRoot)
	if err := http.ListenAndServe(`:8888`, nil); err != nil {
		panic(err)
	}
}
