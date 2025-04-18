package constants

type Module string

type Service string

// Server is the type returned by a classifier server (REST, gRPC)
type Server string

const (
	// REST server
	REST Server = "rest"
	// GRPC server
	GRPC Server = "grpc"
)

func (s Server) String() string {
	return string(s)
}

// Worker is the type returned by a classifier worker
type Worker string

const (
	// RabbitMQ worker
	RabbitMQ Worker = "rabbit-mq"
	// Scheduler worker
	Scheduler Worker = "scheduler"
	// NSQ worker
	NSQ Worker = "nsq"
	// Kafka worker
	Kafka Worker = "kafka"
)

func (w Worker) String() string {
	return string(w)
}
