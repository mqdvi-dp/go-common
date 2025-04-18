package tracer

import (
	"context"
	"fmt"

	sentryexp "github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// initSentry start Sentry
func initSentry(opts Option) (*tracesdk.TracerProvider, Platform) {
	err := sentryexp.Init(
		sentryexp.ClientOptions{
			Dsn:              opts.DSN,
			EnableTracing:    true,
			TracesSampleRate: 1.0,
			Debug:            true,
		},
	)

	if err != nil {
		panic(err)
	}
	// init tracer with otel (open elementary)
	st := tracesdk.NewTracerProvider(
		tracesdk.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)

	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())
	return st, &sentryPlatform{}
}

type sentryPlatform struct{}

func (j *sentryPlatform) Start(ctx context.Context, operationName string) Tracer {
	var span trace.Span
	ctx, span = otel.Tracer(sentry).Start(ctx, operationName)
	if span != nil {
		span = trace.SpanFromContext(ctx)
		ctx = trace.ContextWithSpan(ctx, span)
	}

	return &sentryTracer{span: span, ctx: ctx}
}

func (j *sentryPlatform) GetTraceId(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).TraceID().String()
}

func (j *sentryPlatform) GetSpanId(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).SpanID().String()
}

func (j *sentryPlatform) SetError(ctx context.Context, file string, err error) {
	t := trace.SpanFromContext(ctx)

	t.RecordError(err, trace.WithStackTrace(true))
	t.SetStatus(codes.Error, err.Error())
	t.SetAttributes(attribute.String("error.message", err.Error()))

	logger.Log.ErrorWithFilename(ctx, file, err)
}

func (j *sentryPlatform) Log(ctx context.Context, file, key string, args ...interface{}) {
	t := trace.SpanFromContext(ctx)

	t.SetAttributes(attribute.String(key, fmt.Sprint(args...)))
	if _, ok := keyNotPrinted[key]; ok {
		return
	}

	logger.Log.PrintWithFilename(ctx, file, args...)
}

func (j *sentryPlatform) Debug(ctx context.Context, file, key string, args ...interface{}) {
	if !env.GetBool("DEBUG", false) {
		return
	}

	t := trace.SpanFromContext(ctx)

	t.SetAttributes(attribute.String(key, fmt.Sprint(args...)))
	if _, ok := keyNotPrinted[key]; ok {
		return
	}
	logger.Log.DebugWithFilename(ctx, file, args...)
}

type sentryTracer struct {
	ctx  context.Context
	span trace.Span
	tags map[string]interface{}
}

// Context get active context
func (j *sentryTracer) Context() context.Context {
	return j.ctx
}

// NewContext to continue tracer with new context
func (j *sentryTracer) NewContext() context.Context {
	return trace.ContextWithSpan(context.Background(), j.span)
}

// Tags get tags in tracer span
func (j *sentryTracer) Tags() map[string]interface{} {
	return j.tags
}

// SetTag set tags in tracer span
func (j *sentryTracer) SetTag(key string, value interface{}) {
	if j.span == nil {
		return
	}

	if j.tags == nil {
		j.tags = make(map[string]interface{})
	}

	j.tags[key] = value
}

// SetError set error in span
func (j *sentryTracer) SetError(err error) {
	if j.span == nil || err == nil {
		return
	}

	j.span.RecordError(err, trace.WithStackTrace(true))
	j.span.SetStatus(codes.Error, err.Error())
	j.span.SetAttributes(attribute.String("error.message", err.Error()))
}

// Log log data
func (j *sentryTracer) Log(key string, args ...interface{}) {
	if j.span == nil {
		return
	}

	j.span.SetAttributes(attribute.String(key, fmt.Sprint(args...)))
}

// Debug debug data
func (j *sentryTracer) Debug(key string, args ...interface{}) {
	if !env.GetBool("DEBUG", false) {
		return
	}

	if j.span == nil {
		return
	}

	j.span.SetAttributes(attribute.String(key, fmt.Sprint(args...)))
}

// Finish trace, must in deferred function
func (j *sentryTracer) Finish() {
	if j.span == nil {
		return
	}

	if j.tags != nil || len(j.tags) > 0 {
		for key, value := range j.tags {
			j.span.SetAttributes(attribute.String(key, toValue(value)))
		}
	}

	j.span.End()
}
