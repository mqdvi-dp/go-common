package factory

import (
	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/config"
	"github.com/mqdvi-dp/go-common/constants"
)

// ServiceFactory factory for service
type ServiceFactory interface {
	// GetConfig return an instance of configuration
	GetConfig() config.Config

	// GetDependencies return all dependencies
	GetDependencies() abstract.Dependency

	// GetApplications return all applications (server and/or worker)
	GetApplications() map[string]AppServerFactory

	// RESTHandler return interface of rest-api handler
	RESTHandler() abstract.RESTHandler

	// GRPCHandler return interface of grpc handler
	GRPCHandler() abstract.GRPCHandler

	// WorkerHandler return all interface of worker handler by types.WorkerHandlerGroup
	WorkerHandler(worker constants.Worker) abstract.WorkerHandler

	// Name print service name
	Name() string
}
