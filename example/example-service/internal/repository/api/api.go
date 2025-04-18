package api

import (
	"time"

	"github.com/mqdvi-dp/go-common/request"
	"github.com/mqdvi-dp/go-common/zone"
)

type apiRepository struct {
	client request.ApiClient
	tz     *time.Location
}

func New() *apiRepository {
	return &apiRepository{
		client: request.NewApiClient(),
		tz:     zone.TzJakarta(),
	}
}
