package repository

import (
	"context"

	"github.com/mqdvi-dp/go-common/example/example-service/internal/model"
)

type ApiRepository interface {
	GetFaker(context.Context) (model.ResponseFaker, error)
}

type CacheRepository interface {
	GetFaker(ctx context.Context) (resp model.ResponseFaker, err error)

	SetFaker(ctx context.Context, data model.ResponseFaker) (err error)
}
