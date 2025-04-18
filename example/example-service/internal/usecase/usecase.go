package usecase

import (
	"context"

	"github.com/mqdvi-dp/go-common/example/example-service/internal/model"
	"github.com/mqdvi-dp/go-common/example/example-service/internal/repository"
)

type usecaseInstance struct {
	api   repository.ApiRepository
	cache repository.CacheRepository
}

func New(a repository.ApiRepository, c repository.CacheRepository) *usecaseInstance {
	return &usecaseInstance{
		api:   a,
		cache: c,
	}
}

type Usecase interface {
	GetFaker(ctx context.Context) (resp model.ResponseFaker, err error)
}
