package types

import "time"

// WorkerHandlerFunc handling worker with custom context
type WorkerHandlerFunc func(*EventContext) error

// WorkerHandlerOptionFunc option handler function
type WorkerHandlerOptionFunc func(*WorkerHandler)

// WorkerHandler types
type WorkerHandler struct {
	Pattern      string
	Topic        string
	Queue        string
	ExchangeName string
	HandlerFunc  WorkerHandlerFunc
	DisableTrace bool
	Channel      string
	AutoACK      bool
	MaxRetry     int
	RetryBackoff time.Duration
}

// WorkerHandlerGroup group of worker handlers by pattern
type WorkerHandlerGroup struct {
	Handlers []WorkerHandler
}

// Add method from WorkerHandlerGroup, patternRoute can contain unique topic name, key or task name
func (whg *WorkerHandlerGroup) Add(handlerFunc WorkerHandlerFunc, opts ...WorkerHandlerOptionFunc) {
	h := WorkerHandler{HandlerFunc: handlerFunc, AutoACK: true, MaxRetry: 10, RetryBackoff: 3 * time.Second}

	for _, opt := range opts {
		opt(&h)
	}
	whg.Handlers = append(whg.Handlers, h)
}

// WorkerHandlerOptionDisableTrace set disable trace
func WorkerHandlerOptionDisableTrace() WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.DisableTrace = true
	}
}

// WorkerHandlerOptionChannel set channel
func WorkerHandlerOptionChannel(channel string) WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.Channel = channel
	}
}

// WorkerHandlerOptionTopic set topic
func WorkerHandlerOptionTopic(topic string) WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.Topic = topic
	}
}

// WorkerHandlerOptionQueue set queue
func WorkerHandlerOptionQueue(queue string) WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.Queue = queue
	}
}

// WorkerHandlerOptionExchangeName set exchange name
func WorkerHandlerOptionExchangeName(exchangeName string) WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.ExchangeName = exchangeName
	}
}

// WorkerHandlerOptionAutoACK set auto ack
func WorkerHandlerOptionAutoACK(autoAck bool) WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.AutoACK = autoAck
	}
}

// WorkerHandlerOptionAddHandlers add after handlers execute after main handler
func WorkerHandlerOptionAddHandlers(handlerFuncs WorkerHandlerFunc) WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.HandlerFunc = handlerFuncs
	}
}

// WorkerHandlerOptionPattern set pattern of any worker
func WorkerHandlerOptionPattern(pattern string) WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.Pattern = pattern
	}
}

// WorkerHandlerOptionMaxRetry set max retry
func WorkerHandlerOptionMaxRetry(maxRetry int) WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.MaxRetry = maxRetry
	}
}

// WorkerHandlerOptionRetryBackoff set retry backoff
func WorkerHandlerOptionRetryBackoff(retryBackoff time.Duration) WorkerHandlerOptionFunc {
	return func(wh *WorkerHandler) {
		wh.RetryBackoff = retryBackoff
	}
}
