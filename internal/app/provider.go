package app

import (
	"github.com/blagorodov/go-shortener/internal/repository"
	"github.com/blagorodov/go-shortener/internal/service"
)

type Provider struct {
	Repository repository.Repository
	Service    service.Service
}

func newProvider() *Provider {
	return &Provider{}
}
