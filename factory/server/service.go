package server

import (
	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/config"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/factory"
	"github.com/mqdvi-dp/go-common/factory/server/cron"
	"github.com/mqdvi-dp/go-common/factory/server/kafka"
	"github.com/mqdvi-dp/go-common/factory/server/nsq"
	"github.com/mqdvi-dp/go-common/factory/server/rest"
	"github.com/mqdvi-dp/go-common/factory/server/rmq"
	"github.com/mqdvi-dp/go-common/factory/server/rpc"
)

// ServiceFunc setter to set service instance
type ServiceFunc func(*service)

// service instance
type service struct {
	config               config.Config
	dependencies         abstract.Dependency
	restHandler          abstract.RESTHandler
	restHandlerOptions   []rest.OptionFunc
	grpcHandler          abstract.GRPCHandler
	grpcHandlerOptions   []rpc.OptionFunc
	workerHandler        map[constants.Worker]abstract.WorkerHandler
	workerHandlerOptions map[constants.Worker]interface{}
	applications         map[string]factory.AppServerFactory
}

// SetConfiguration sets the configuration
func SetConfiguration(cfg config.Config) ServiceFunc {
	return func(s *service) {
		s.config = cfg
	}
}

// SetDependencies sets the services dependencies
func SetDependencies(dpnds abstract.Dependency) ServiceFunc {
	return func(s *service) {
		s.dependencies = dpnds
	}
}

// SetRestHandler sets handler for rest-api
func SetRestHandler(restHandler abstract.RESTHandler, opts ...rest.OptionFunc) ServiceFunc {
	return func(s *service) {
		s.restHandler = restHandler
		s.restHandlerOptions = opts
	}
}

// SetGrpcHandler sets handler for gRPC
func SetGrpcHandler(grpcHandler abstract.GRPCHandler, opts ...rpc.OptionFunc) ServiceFunc {
	return func(s *service) {
		s.grpcHandler = grpcHandler
		s.grpcHandlerOptions = opts
	}
}

// SetWorkerHandler sets handler for rest-api
func SetWorkerHandler(worker constants.Worker, workerHandler abstract.WorkerHandler, workerOptions ...interface{}) ServiceFunc {
	return func(s *service) {
		if s.workerHandler == nil || len(s.workerHandler) < 1 {
			s.workerHandler = make(map[constants.Worker]abstract.WorkerHandler)
		}

		s.workerHandler[worker] = workerHandler

		if s.workerHandlerOptions == nil || len(s.workerHandlerOptions) < 1 {
			s.workerHandlerOptions = make(map[constants.Worker]interface{})
		}

		s.workerHandlerOptions[worker] = workerOptions
	}
}

// NewApplicationService creates a new service instance
func NewApplicationService(funcs ...ServiceFunc) factory.ServiceFactory {
	svc := &service{}

	// register factories
	for _, f := range funcs {
		f(svc)
	}

	return svc
}

func (s *service) GetConfig() config.Config {
	return s.config
}

func (s *service) GetDependencies() abstract.Dependency {
	return s.dependencies
}

func (s *service) Name() string {
	return s.config.GetServiceName()
}

func (s *service) RESTHandler() abstract.RESTHandler {
	return s.restHandler
}

func (s *service) GRPCHandler() abstract.GRPCHandler {
	return s.grpcHandler
}

func (s *service) WorkerHandler(worker constants.Worker) abstract.WorkerHandler {
	return s.workerHandler[worker]
}

func (s *service) GetApplications() map[string]factory.AppServerFactory {
	// initiate when map not yet declare
	// handling error nil pointer reference
	if len(s.applications) < 1 || s.applications == nil {
		s.applications = make(map[string]factory.AppServerFactory)
	}

	// when rest handler is empty, create default application for rest handler
	if s.restHandler == nil {
		s.restHandler = defaultRestHandler()
	}

	// check is rest already registered
	if _, ok := s.applications[constants.REST.String()]; !ok {
		// initialized application rest server
		if s.restHandler != nil {
			s.applications[constants.REST.String()] = rest.New(s, s.restHandlerOptions...)
		}
	}

	// check is grpc already registered
	if _, ok := s.applications[constants.GRPC.String()]; !ok {
		// initialized application grpc server
		if s.grpcHandler != nil {
			s.applications[constants.GRPC.String()] = rpc.New(s, s.grpcHandlerOptions...)
		}
	}

	// is have worker handler for rabbit-mq?
	if s.workerHandler[constants.RabbitMQ] != nil {
		// check is rabbit-mq already registered
		if _, ok := s.applications[constants.RabbitMQ.String()]; !ok {
			if s.workerHandler[constants.RabbitMQ] != nil {
				var rmqOptions []rmq.OptionFunc
				if val, ok := s.workerHandlerOptions[constants.RabbitMQ]; ok {
					if intfs, ok := val.([]interface{}); ok {
						for _, intf := range intfs {
							if opt, ok := intf.(rmq.OptionFunc); ok {
								rmqOptions = append(rmqOptions, opt)
							}
						}
					}
				}

				// initialized application rabbit-mq consumer and publisher
				s.applications[constants.RabbitMQ.String()] = rmq.NewWorker(s, rmqOptions...)
			}
		}
	}

	// is have worker handler for nsq?
	if s.workerHandler[constants.NSQ] != nil {
		// check is nsq already registered
		if _, ok := s.applications[constants.NSQ.String()]; !ok {
			if s.workerHandler[constants.NSQ] != nil {
				var nsqOptions []nsq.OptionFunc
				if val, ok := s.workerHandlerOptions[constants.NSQ]; ok {
					if intfs, ok := val.([]interface{}); ok {
						for _, intf := range intfs {
							if opt, ok := intf.(nsq.OptionFunc); ok {
								nsqOptions = append(nsqOptions, opt)
							}
						}
					}
				}

				// initialized application nsq consumer
				s.applications[constants.NSQ.String()] = nsq.New(s, nsqOptions...)
			}
		}
	}

	// is have worker handler for scheduller?
	if s.workerHandler[constants.Scheduler] != nil {
		// check is scheduler already registered
		if _, ok := s.applications[constants.Scheduler.String()]; !ok {
			if s.workerHandler[constants.Scheduler] != nil {
				var cronOptions []cron.OptionFunc
				if val, ok := s.workerHandlerOptions[constants.Scheduler]; ok {
					if intfs, ok := val.([]interface{}); ok {
						for _, intf := range intfs {
							if opt, ok := intf.(cron.OptionFunc); ok {
								cronOptions = append(cronOptions, opt)
							}
						}
					}
				}

				// initialized application scheduler consumer and publisher
				s.applications[constants.Scheduler.String()] = cron.NewWorker(s, cronOptions...)
			}
		}
	}

	// is have worker handler for kafka?
	if s.workerHandler[constants.Kafka] != nil {
		// check is kafka already registered
		if _, ok := s.applications[constants.Kafka.String()]; !ok {
			if s.workerHandler[constants.Kafka] != nil {
				var kafkaOptions []kafka.OptionFunc
				if val, ok := s.workerHandlerOptions[constants.Kafka]; ok {
					if intfs, ok := val.([]interface{}); ok {
						for _, intf := range intfs {
							if opt, ok := intf.(kafka.OptionFunc); ok {
								kafkaOptions = append(kafkaOptions, opt)
							}
						}
					}
				}

				// initialized application kafka consumer
				s.applications[constants.Kafka.String()] = kafka.New(s, kafkaOptions...)
			}
		}
	}

	return s.applications
}
