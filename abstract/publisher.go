package abstract

import (
	"context"

	"github.com/mqdvi-dp/go-common/types"
)

// Publisher abstraction
type Publisher interface {
	PublishMessage(ctx context.Context, req *types.PublisherArgument) error
	PublishMessages(ctx context.Context, req []*types.PublisherArgument) error
}
