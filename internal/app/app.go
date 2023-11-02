package app

import (
	"context"
	"github.com/blagorodov/go-shortener/internal/config"
	"github.com/blagorodov/go-shortener/internal/logger"
	"github.com/blagorodov/go-shortener/internal/repository"
	"github.com/blagorodov/go-shortener/internal/repository/db"
	"github.com/blagorodov/go-shortener/internal/repository/file"
	"github.com/blagorodov/go-shortener/internal/repository/memory"
	"github.com/blagorodov/go-shortener/internal/server"
	"github.com/blagorodov/go-shortener/internal/service/shortener"
)

type App struct {
	Provider *Provider
	Server   *server.Server
	Ctx      context.Context
}

func Create(ctx context.Context) (*App, error) {
	var rp repository.Repository
	var err error
	logger.Init()
	switch {
	case config.Options.DBDataSource != "":
		rp, err = db.NewRepository(ctx)
	case config.Options.URLDBPath != "":
		rp, err = file.NewRepository(ctx)
	default:
		rp, err = memory.NewRepository(ctx)
	}
	if err != nil {
		return nil, err
	}

	a := &App{}
	a.Ctx = ctx
	a.Provider = &Provider{
		Repository: rp,
		Service:    shortener.NewService(rp),
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
