package factory

import "context"

// AppServerFactory factory for server and/or worker abstraction
type AppServerFactory interface {
	// Serve for running server or worker
	Serve()
	// Shutdown stop the server or worker
	Shutdown(ctx context.Context)
	// Name print name of server or worker
	Name() string
}
