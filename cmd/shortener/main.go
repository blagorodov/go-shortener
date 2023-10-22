package main

import (
	"context"
	app2 "github.com/blagorodov/go-shortener/internal/app"
)

func main() {
	app, err := app2.Create(context.Background())
	if err != nil {
		panic(err)
	}
	defer app.Destroy()

	app.Run()
}
