package server

import (
	"github.com/blagorodov/go-shortener/internal/app/handlers"
	"net/http"
)

// Start Запуск сервера
func Start() {
	http.HandleFunc(`/`, handlers.Root)
	createServer()
}

func createServer() {
	if err := http.ListenAndServe(`:8080`, nil); err != nil {
		panic(err)
	}
}
