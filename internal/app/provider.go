package app

import (
	"github.com/blagorodov/go-shortener/internal/repository"
	"github.com/blagorodov/go-shortener/internal/service"
)

type provider struct {
	repository repository.Repository
	service    service.Service
}

func newProvider() *provider {
	return &provider{}
}
