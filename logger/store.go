package logger

import (
	"context"
	"time"
)

// Store the outgoing request into the context
func (o OutgoingLog) Store(ctx context.Context) {
	var og = make([]OutgoingLog, 0)

	value, ok := extract(ctx)
	if !ok {
		return
	}

	// get key and delete from context
	tmp, ok := value.LoadAndDelete(_OutgoingLog)
	if ok {
		og = tmp.([]OutgoingLog)
	}

	og = append(og, o)
	// set key into context with new value
	value.Set(_OutgoingLog, og)
}

// Store the database logging in the context
func (d Database) Store(ctx context.Context) {
	var db = make([]Database, 0)

	value, ok := extract(ctx)
	if !ok {
		return
	}

	// get key and delete from context
	tmp, ok := value.LoadAndDelete(_Database)
	if ok {
		db = tmp.([]Database)
	}

	d.ExecutionTime = time.Since(d.startTime).Seconds()
	db = append(db, d)
	// set key into context with new value
	value.Set(_Database, db)
}

// SetUsername sets the username into the context
func SetUsername(ctx context.Context, username string) {
	Log.Debugf(ctx, "username: %s", username)
	value, ok := extract(ctx)
	if !ok {
		return
	}

	value.Set(_Username, username)
}

// SetErrorMessage sets the error message into context
func SetErrorMessge(ctx context.Context, err error) {
	if err == nil {
		return
	}

	value, ok := extract(ctx)
	if !ok {
		return
	}

	value.Set(_ErrorMessage, err.Error())
}
