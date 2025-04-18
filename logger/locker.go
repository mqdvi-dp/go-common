package logger

import (
	"context"
)

// NewLocker creates a new instance of Locker that will be used for locking and unlocking logging messages
func NewLocker(ctxs ...context.Context) *Locker {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	if len(ctxs) > 0 {
		ctx, cancel = context.WithCancel(ctxs[0])
	}

	return &Locker{ctx: ctx, cancel: cancel}
}

// Values create contract store value into context
type Values interface {
	Set(key Flags, value interface{})
	Load(key Flags) (interface{}, bool)
	LoadAndDelete(key Flags) (interface{}, bool)
	Delete(key Flags)
	Cleanup()
}

// Set value to keys
func (l *Locker) Set(key Flags, value interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.data.Store(key, value)
}

// Delete value from sync
func (l *Locker) Delete(key Flags) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.data.Delete(key)
}

// Load value from key
func (l *Locker) Load(key Flags) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.data.Load(key)
}

// LoadAndDelete from key
func (l *Locker) LoadAndDelete(key Flags) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.data.LoadAndDelete(key)
}

func (l *Locker) Cleanup() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.cancel()
}

func extract(ctx context.Context) (Values, bool) {
	var (
		lock = new(Locker)
		ok   bool
	)

	if ctx == nil {
		return lock, false
	}

	lock, ok = ctx.Value(LogKey).(*Locker)
	return lock, ok
}
