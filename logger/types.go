package logger

import (
	"context"
	"sync"
	"time"
)

type (
	// Key of context
	Key int
	// Flags is key for store context
	Flags string
	// HandlerType is type for data logging
	HandlerType string
	// DatabaseType is type for database logging
	DatabaseType string
)

const (
	// LogKey is key context for handler type
	LogKey Key = 31

	// Http is type for logging Http Request
	Http HandlerType = "http"
	// GRPC is type for logging GRPC Request
	GRPC HandlerType = "grpc"
	// NSQ is type for logging NSQ Consumer
	NSQ HandlerType = "nsq_consumer"
	// RabbitMQ is type for logging RabbitMQ Consumer
	RabbitMQ HandlerType = "rabbitmq_consumer"
	// Kafka is type for logging Kafka Consumer
	Kafka HandlerType = "kafka_consumer"
	// Scheduler is type for logging Scheduler (cron job)
	Scheduler HandlerType = "scheduler"

	// Sql is database logging type for SQL
	Sql DatabaseType = "sql"
	// SqlTx is database logging type for SQLTx
	SqlTx DatabaseType = "sqlTx"
	// Redis is database logging type for Redis
	Redis DatabaseType = "redis"

	// _Username is key of logging field struct
	_Username Flags = "Username"
	// _OutgoingLog is key of logging field struct
	_OutgoingLog Flags = "OutgoingLog"
	// _LogMessage is key of logging field struct
	_LogMessage Flags = "LogMessage"
	// _Database is key of logging field struct
	_Database Flags = "Database"
	// _RequestId is key of logging field struct
	_RequestId Flags = "RequestId"
	// _ErrorMessage is key of logging field struct
	_ErrorMessage Flags = "ErrorMessage"
)

// Logger standard data model for logging
type Logger struct {
	HandlerType   HandlerType   `json:"handler_type"`
	Ip            string        `json:"ip"`
	DeviceId      string        `json:"device_id"`
	Timezone      string        `json:"timezone"`
	Lat           float64       `json:"lat"`
	Lng           float64       `json:"lng"`
	Source        string        `json:"source"`
	StartTime     string        `json:"start_time"`
	RequestId     string        `json:"request_id"`
	Username      string        `json:"username"`
	Service       string        `json:"service"`
	Endpoint      string        `json:"endpoint"`
	RequestHeader string        `json:"request_header"`
	RequestBody   string        `json:"request_body"`
	ResponseBody  string        `json:"response_body"`
	StatusCode    int           `json:"status_code"`
	ErrorMessage  string        `json:"error_message"`
	ExecutionTime float64       `json:"execution_time"`
	LogMessage    []LogMessage  `json:"debugging"`
	OutgoingLog   []OutgoingLog `json:"outgoing_log"`
	Database      []Database    `json:"database"`
}

// OutgoingLog represents the state of a outgoing log request.
type OutgoingLog struct {
	StartTime     string  `json:"start_time"`
	TargetService string  `json:"target_service"`
	URL           string  `json:"url"`
	RequestHeader string  `json:"request_header"`
	RequestBody   string  `json:"request_body"`
	ResponseBody  string  `json:"response_body"`
	StatusCode    int     `json:"status_code"`
	ExecutionTime float64 `json:"execution_time"`
}

// LogMessage represents the state of a debugging log.
type LogMessage struct {
	File    string `json:"file"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

// Database represents the state of a database log
type Database struct {
	Type          DatabaseType `json:"type"`
	Query         string       `json:"query"`
	Arguments     []string     `json:"arguments"`
	startTime     time.Time    `json:"-"`
	ExecutionTime float64      `json:"execution_time"`
}

// Locker is container data
type Locker struct {
	ctx    context.Context
	cancel context.CancelFunc
	data   sync.Map
	mutex  sync.Mutex
}
