package main

import (
	"context"
	"time"

	"github.com/mqdvi-dp/go-common/tracer"
)

func main() {
	tracer.New(
		"sentry-example",
		tracer.SetPlatformType(tracer.Sentry),
		tracer.SetDSN("https://98b83ee7c2219ab2083a38d41e1c3041@o4506048835223552.ingest.sentry.io/4506048836927488"),
	)

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log(ctx)
}

func log(ctx context.Context) {
	trace := tracer.StartTrace(ctx, "Usecase:Login")
	defer trace.Finish()
}
