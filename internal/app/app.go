package app

import (
	"context"
	"github.com/blagorodov/go-shortener/internal/logger"
	"github.com/blagorodov/go-shortener/internal/repository/file"
	"github.com/blagorodov/go-shortener/internal/server"
	"github.com/blagorodov/go-shortener/internal/service/shortener"
)

type App struct {
	Provider *Provider
	Server   *server.Server
	Ctx      context.Context
}

func Create(ctx context.Context) (*App, error) {
	logger.Init()

	a := &App{
		Ctx: ctx,
	}
	fileRepository, err := file.NewRepository(ctx)
	if err != nil {
		return nil, err
	}
	a.Provider = &Provider{
		Repository: fileRepository,
		Service:    shortener.NewService(fileRepository),
	}
	a.Server = server.NewServer(ctx, a.Provider.Service)
	a.Server.Start()

	return a, nil
}

func (a *App) Run() {
	a.Server.Start()
}

func (a *App) Destroy() error {
	return a.Provider.Repository.Destroy()
}
