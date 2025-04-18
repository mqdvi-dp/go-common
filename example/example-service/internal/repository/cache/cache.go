package cache

import "github.com/mqdvi-dp/go-common/config/database/rdc"

const (
	prefixKeyFaker = "faker"
)

type cacheRepository struct {
	client rdc.Rdc
}

func New(c rdc.Rdc) *cacheRepository {
	return &cacheRepository{client: c}
}
