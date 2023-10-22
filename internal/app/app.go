package app

import (
	"context"
	"github.com/blagorodov/go-shortener/internal/logger"
	"github.com/blagorodov/go-shortener/internal/repository/file"
	"github.com/blagorodov/go-shortener/internal/server"
	"github.com/blagorodov/go-shortener/internal/service/shortener"
)

type App struct {
	provider *provider
	server   *server.Server
	ctx      context.Context
}

func Create(ctx context.Context) (*App, error) {
	logger.Init()

	a := &App{
		ctx: ctx,
	}
	fileRepository, err := file.NewRepository(ctx)
	if err != nil {
		return nil, err
	}
	a.provider = &provider{
		repository: fileRepository,
		service:    shortener.NewService(fileRepository),
	}
	a.server = server.NewServer(ctx, a.provider.service)
	a.server.Start()

	return a, nil
}

func (a *App) Run() {
	a.server.Start()
}

func (a *App) Destroy() error {
	return a.provider.repository.Destroy()
}
