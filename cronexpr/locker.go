package cronexpr

import (
	"context"
	"strconv"

	"github.com/mqdvi-dp/go-common/config/database/rdc"
)

type redisLocker struct {
	pool rdc.Rdc
}

type NoopLocker struct{}

// Locker abstraction, lock concurrent process
type Locker interface {
	IsLocked(key string) bool
	HasBeenLocked(key string) bool
	Unlock(key string)
	Reset(key string)
}

// NewRedisLocker constructor
func NewRedisLocker(pool rdc.Rdc) Locker {
	return &redisLocker{pool: pool}
}

func (r *redisLocker) IsLocked(key string) bool {
	incr, _ := r.pool.Incr(context.Background(), key)

	return incr > 1
}

func (r *redisLocker) HasBeenLocked(key string) bool {
	result, _ := r.pool.Get(context.Background(), key)
	incr, _ := strconv.Atoi(result)

	return incr > 0
}

func (r *redisLocker) Unlock(key string) {
	_ = r.pool.Del(context.Background(), key)
}

func (r *redisLocker) Reset(key string) {
	_ = r.pool.Del(context.Background(), key)
}

// IsLocked method
func (NoopLocker) IsLocked(key string) bool {
	return false
}

// HasBeenLocked method
func (NoopLocker) HasBeenLocked(key string) bool {
	return false
}

// Unlock method
func (NoopLocker) Unlock(key string) {
}

// Reset method
func (NoopLocker) Reset(key string) {
}
