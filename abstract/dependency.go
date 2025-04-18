package abstract

import "github.com/mqdvi-dp/go-common/constants"

type depend struct {
	mw      Middleware
	sqlDb   map[Instance]SQLDatabase
	redisDb RedisDatabase
	brokers map[constants.Worker]Broker
}

type Option func(*depend)

// SetMiddleware option func to set Middleware dependencies
func SetMiddleware(mw Middleware) Option {
	return func(d *depend) {
		d.mw = mw
	}
}

// SetSQLDatabase option func to set SQLDatabase dependencies
func SetSQLDatabase(instance Instance, db SQLDatabase) Option {
	return func(d *depend) {
		if d.sqlDb == nil {
			d.sqlDb = make(map[Instance]SQLDatabase)
		}

		d.sqlDb[instance] = db
	}
}

// SetRedisDatabase option func to set RedisDatabase dependencies
func SetRedisDatabase(db RedisDatabase) Option {
	return func(d *depend) {
		d.redisDb = db
	}
}

// SetBrokers option func to set Broker dependencies
func SetBrokers(brokers map[constants.Worker]Broker) Option {
	return func(d *depend) {
		d.brokers = brokers
	}
}

type Dependency interface {
	// GetMiddleware get current middleware
	GetMiddleware() Middleware
	// SetMiddleware set custom middleware
	SetMiddleware(mw Middleware)

	// GetSQLDatabase database sql dependencies
	GetSQLDatabase(instance Instance) SQLDatabase
	// GetRedisDatabase database redis
	GetRedisDatabase() RedisDatabase

	// GetBroker get broker based on parameters
	GetBroker(constants.Worker) Broker
	// AddBroker add custom broker into variable brokers
	AddBroker(worker constants.Worker, broker Broker)
}

var stdDeps = new(depend)

func New(opts ...Option) Dependency {
	for _, o := range opts {
		o(stdDeps)
	}

	return stdDeps
}

func (d *depend) GetMiddleware() Middleware {
	return d.mw
}

func (d *depend) SetMiddleware(mw Middleware) {
	d.mw = mw
}

func (d *depend) GetSQLDatabase(instance Instance) SQLDatabase {
	return d.sqlDb[instance]
}

func (d *depend) GetRedisDatabase() RedisDatabase {
	return d.redisDb
}

func (d *depend) GetBroker(worker constants.Worker) Broker {
	return d.brokers[worker]
}

func (d *depend) AddBroker(worker constants.Worker, broker Broker) {
	if d.brokers == nil {
		d.brokers = make(map[constants.Worker]Broker)
	}

	d.brokers[worker] = broker
}
