package tracer

import (
	"context"
	"sync"
)

type PlatformType string

const (
	// Jaeger tracer platform
	Jaeger PlatformType = "jaeger"
	// Sentry tracer platform
	Sentry PlatformType = "sentry"
)

var (
	once         sync.Once
	activeTracer Platform
)

type Platform interface {
	Start(ctx context.Context, operationName string) Tracer
	GetTraceId(ctx context.Context) string
	GetSpanId(ctx context.Context) string
	SetError(ctx context.Context, file string, err error)
	Log(ctx context.Context, file, key string, args ...interface{})
	Debug(ctx context.Context, file, key string, args ...interface{})
}

// SetTracerPlatformType function for set tracer platform
func SetTracerPlatformType(t Platform) {
	once.Do(
		func() {
			activeTracer = t
		},
	)
}
