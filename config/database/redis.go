package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/config/database/rdc"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/redis/go-redis/v9"
)

type RedisFuncOption func(*redisOption)

type redisOption struct {
	serviceName       string
	address           string
	dbNumber          int
	poolSize          int
	minIdleConnection int
	maxIdleConnection int
	maxIdleTimeout    time.Duration
}

func defaultRedisOption() redisOption {
	return redisOption{
		address:           env.GetString("REDIS_MACHINES", "redis://localhost:6379"),
		dbNumber:          env.GetInt("REDIS_DB_NUMBER", 0),
		poolSize:          env.GetInt("REDIS_POOL_SIZE", 20),
		minIdleConnection: env.GetInt("REDIS_MIN_IDLE_CONNECTION", 0),
		maxIdleConnection: env.GetInt("REDIS_MAX_IDLE_CONNECTION", 0),
		maxIdleTimeout:    env.GetDuration("REDIS_MAX_IDLE_TIMEOUT", time.Duration(5)*time.Second),
	}
}

type redisInstance struct {
	db rdc.Rdc
}

func (r *redisInstance) Client() rdc.Rdc {
	return r.db
}

func (r *redisInstance) Disconnect(ctx context.Context) error {
	logger.RedBold("redis: disconnecting...")
	defer fmt.Printf("\x1b[31;1mRedis Disconnecting:\x1b[0m \x1b[32;1mSUCCESS\x1b[0m\n")

	return r.db.Close()
}

// NewRedisConnection returns a Redis connection
func NewRedisConnection(opts ...RedisFuncOption) abstract.RedisDatabase {
	logger.YellowItalic("Load redis connection...")
	// redis custom option
	opt := defaultRedisOption()
	// merge with parameters
	for _, o := range opts {
		o(&opt)
	}

	var redisDriver = env.GetString("REDIS_DRIVER", "redis")
	var client redis.UniversalClient
	switch redisDriver {
	case "redis-sentinel":
		// initialize connection
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			ClientName:      opt.serviceName,
			DB:              opt.dbNumber,
			MasterName:      "mymaster",
			SentinelAddrs:   strings.Split(opt.address, ","),
			PoolSize:        opt.poolSize,
			MinIdleConns:    opt.minIdleConnection,
			MaxIdleConns:    opt.maxIdleConnection,
			ConnMaxIdleTime: opt.maxIdleTimeout,
			RouteByLatency:  false,
			RouteRandomly:   false,
		})
	case "redis":
		client = redis.NewClient(&redis.Options{
			ClientName:      opt.serviceName,
			Addr:            opt.address,
			DB:              opt.dbNumber,
			PoolSize:        opt.poolSize,
			MinIdleConns:    opt.minIdleConnection,
			MaxIdleConns:    opt.maxIdleConnection,
			ConnMaxIdleTime: opt.maxIdleTimeout,
		})
	case "redis-cluster":
		// initialize connection
		client = redis.NewClusterClient(&redis.ClusterOptions{
			ClientName:      opt.serviceName,
			Addrs:           strings.Split(opt.address, ","),
			PoolSize:        opt.poolSize,
			MinIdleConns:    opt.minIdleConnection,
			MaxIdleConns:    opt.maxIdleConnection,
			ConnMaxIdleTime: opt.maxIdleTimeout,
		})
	}

	// creates context for timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	logger.GreenItalic("redis connected!")
	return &redisInstance{db: &rdc.Db{DB: client}}
}

// SetRedisServiceName sets the service name
func SetRedisServiceName(serviceName string) RedisFuncOption {
	return func(ro *redisOption) {
		ro.serviceName = serviceName
	}
}

// SetRedisAddress sets the redis address
func SetRedisAddress(address []string) RedisFuncOption {
	return func(ro *redisOption) {
		ro.address = strings.Join(address, ",")
	}
}

// SetRedisDbNumber sets the redis database number
func SetRedisDbNumber(dbNumber int) RedisFuncOption {
	return func(ro *redisOption) {
		ro.dbNumber = dbNumber
	}
}

// SetRedisPoolSize sets the redis pool size
func SetRedisPoolSize(poolSize int) RedisFuncOption {
	return func(ro *redisOption) {
		ro.poolSize = poolSize
	}
}

// SetRedisMinIdleConnection sets the minimum idle connection of redis
func SetRedisMinIdleConnection(minIdleConnection int) RedisFuncOption {
	return func(ro *redisOption) {
		ro.minIdleConnection = minIdleConnection
	}
}

// SetRedisMaxIdleConnection sets the maximum idle connection of redis
func SetRedisMaxIdleConnection(maxIdleConnection int) RedisFuncOption {
	return func(ro *redisOption) {
		ro.maxIdleConnection = maxIdleConnection
	}
}

// SetRedisMaxIdleTimeout sets the maximum idle timeout
func SetRedisMaxIdleTimeout(maxIdleTimeout time.Duration) RedisFuncOption {
	return func(ro *redisOption) {
		ro.maxIdleTimeout = maxIdleTimeout
	}
}
