package tracer

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/mqdvi-dp/go-common/logger"
	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

var keyNotPrinted = map[string]bool{
	"query":            true,
	"argument":         true,
	"arguments":        true,
	"args":             true,
	"exec":             true,
	"result":           true,
	"key":              true,
	"expired_duration": true,
}

type noopTracerPlatform struct {
	ctx context.Context
}

type Tracer interface {
	Context() context.Context
	NewContext() context.Context
	Tags() map[string]interface{}
	SetTag(key string, value interface{})
	Log(key string, args ...interface{})
	Debug(key string, args ...interface{})
	SetError(err error)
	Finish()
}

// New initiate tracer with some Options
func New(serviceName string, opts ...OptionFunc) {
	var (
		platform Platform
		tracer   *tracesdk.TracerProvider
	)

	opt := defaultOption
	opt.ServiceName = serviceName
	for _, o := range opts {
		o(&opt)
	}

	// updated: 14 October 2023
	// jaeger exporter is no longer support in opentelemetry package
	// so, platformType should set, if not, we will send the fatal error
	//
	// support package:
	// - sentry
	switch opt.platformType {
	case Jaeger:
		log.Fatalln(fmt.Errorf("jaeger exporters is no longer supported, please use Sentry"))
	case Sentry:
		tracer, platform = initSentry(opt)
	default: // default is jaeger
		log.Fatalln(fmt.Errorf("platform type is not set"))
	}

	if tracer == nil || platform == nil {
		log.Fatalln(fmt.Errorf("tracer or platform is nil"))
	}

	otel.SetTracerProvider(tracer)
	SetTracerPlatformType(platform)
}

// StartTrace starting trace child span from parent span
func StartTrace(ctx context.Context, operationName string) Tracer {
	if activeTracer == nil {
		activeTracer = noopTracerPlatform{ctx: ctx}
	}

	return activeTracer.Start(ctx, operationName)
}

// StartTraceWithContext starting to trace child span from parent span with context
func StartTraceWithContext(ctx context.Context, operationName string) (Tracer, context.Context) {
	t := StartTrace(ctx, operationName)

	return t, t.Context()
}

// GetTraceId get active trace id
func GetTraceId(ctx context.Context) string {
	if activeTracer == nil {
		activeTracer = noopTracerPlatform{ctx: ctx}
	}

	return activeTracer.GetTraceId(ctx)
}

// GetSpanId get current span id
func GetSpanId(ctx context.Context) string {
	if activeTracer == nil {
		activeTracer = noopTracerPlatform{ctx: ctx}
	}

	return activeTracer.GetSpanId(ctx)
}

// SetError set error into active span
func SetError(ctx context.Context, err error) {
	if activeTracer == nil {
		activeTracer = &noopTracerPlatform{ctx: ctx}
	}

	var file string
	// get file and line
	_, fileName, line, _ := runtime.Caller(1)
	// default file
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	activeTracer.SetError(ctx, file, err)
}

// Log set log attributes into active span
func Log(ctx context.Context, key string, args ...interface{}) {
	if activeTracer == nil {
		activeTracer = &noopTracerPlatform{ctx: ctx}
	}

	var file string
	// get file and line
	_, fileName, line, _ := runtime.Caller(1)
	// default file
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	activeTracer.Log(ctx, file, key, args...)
}

// Debug set log-debug attributes into active span
func Debug(ctx context.Context, key string, args ...interface{}) {
	if activeTracer == nil {
		activeTracer = &noopTracerPlatform{ctx: ctx}
	}

	var file string
	// get file and line
	_, fileName, line, _ := runtime.Caller(1)
	// default file
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	activeTracer.Debug(ctx, file, key, args...)
}

func (n noopTracerPlatform) Start(ctx context.Context, operationName string) Tracer {
	return noopTracer{ctx: ctx}
}

func (n noopTracerPlatform) GetTraceId(ctx context.Context) string {
	return ""
}

func (n noopTracerPlatform) GetSpanId(ctx context.Context) string {
	return ""
}

func (n noopTracerPlatform) SetError(ctx context.Context, file string, err error) {
	logger.Log.ErrorWithFilename(ctx, file, err)
}

func (n noopTracerPlatform) Log(ctx context.Context, file, key string, args ...interface{}) {
	logger.Log.PrintWithFilename(ctx, file, args...)
}

func (n noopTracerPlatform) Debug(ctx context.Context, file, key string, args ...interface{}) {
	logger.Log.DebugWithFilename(ctx, file, args...)
}

type noopTracer struct {
	ctx context.Context
}

func (n noopTracer) Context() context.Context {
	return n.ctx
}

func (n noopTracer) NewContext() context.Context {
	return n.ctx
}

func (noopTracer) Tags() map[string]interface{} {
	return map[string]interface{}{}
}

func (noopTracer) SetTag(key string, value interface{}) {}

func (noopTracer) SetError(err error) {
}

func (noopTracer) Log(key string, args ...interface{}) {}

func (noopTracer) Debug(key string, args ...interface{}) {}

func (noopTracer) Finish() {}

func (noopTracer) GetTraceId(ctx context.Context) (u string) {
	return
}

func (noopTracer) GetSpanId(ctx context.Context) (u string) {
	return
}
