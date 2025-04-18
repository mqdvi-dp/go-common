package rdc

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

// validate is a helper function to validate value should be pointers
func validate(val interface{}) error {
	vof := reflect.ValueOf(val)

	if vof.Kind() != reflect.Ptr {
		return fmt.Errorf("destination should be pointer")
	}

	return nil
}

type Db struct {
	DB redis.UniversalClient
}

type Rdc interface {
	// Get returns single data based on key
	Get(ctx context.Context, key string) (string, error)

	// GetStruct returns single data and will marshal into struct based on key
	GetStruct(ctx context.Context, dest interface{}, key string) error

	// Incr value based on key
	Incr(ctx context.Context, key string) (int64, error)

	// Decr value based on key
	Decr(ctx context.Context, key string) (int64, error)

	// Set value into redis with data type string
	// default of expired duration is until end of day
	Set(ctx context.Context, key string, value interface{}, durations ...time.Duration) error

	// Del value from redis
	Del(ctx context.Context, keys ...string) error

	// TTL get time to live of the key
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Expire set the expired time of the key
	Expire(ctx context.Context, key string, duration time.Duration) error

	// Keys get all keys related that pattern
	Keys(ctx context.Context, patternKey string) ([]string, error)

	// GetKeysAndDelete get all keys based on patterns and delete the keys after that
	GetKeysAndDelete(ctx context.Context, patternKey string) (err error)

	// Returns all the members of the set value stored at key.
	DoSMembers(ctx context.Context, key string) (members []string, err error)

	// Returns true if member is a member of the set stored at key, otherwise false
	DoSIsMember(ctx context.Context, key, member string) (exist bool)

	// Add the specified members with array of string data type to the set stored at key
	DoSadd(ctx context.Context, key string, value []string, durations ...time.Duration) (err error)

	// Set value with field
	// default of expired duration is until end of day
	HSet(ctx context.Context, key string, field string, value interface{}, durations ...time.Duration) error

	// HGet returns single value from selected key and field
	HGet(ctx context.Context, key string, field string) (string, error)

	// HGetAll returns all fields under selected key
	HGetAll(ctx context.Context, key string) (map[string]string, error)

	// Close the connection
	Close() error

	// Check the connection
	Ping(ctx context.Context) error
}
