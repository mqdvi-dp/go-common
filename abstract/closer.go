package abstract

import "context"

// Closer abstraction
type Closer interface {
	Disconnect(ctx context.Context) error
}
