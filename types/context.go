package types

import (
	"bytes"
	"context"
)

type EventContext struct {
	ctx        context.Context
	workerType string
	topic      string
	header     map[string]interface{}
	key        string
	err        error
	buff       *bytes.Buffer
}

// NewEventContext event context constructor
func NewEventContext(buff *bytes.Buffer) *EventContext {
	buff.Reset()
	return &EventContext{
		buff: buff,
	}
}

// SetContext setter context
func (e *EventContext) SetContext(ctx context.Context) {
	e.ctx = ctx
}

// SetWorkerType setter worker type
func (e *EventContext) SetWorkerType(wt string) {
	e.workerType = wt
}

// SetTopic setter handler route
func (e *EventContext) SetTopic(h string) {
	e.topic = h
}

// SetHeader setter header
func (e *EventContext) SetHeader(header map[string]interface{}) {
	e.header = header
}

// SetKey setter key context
func (e *EventContext) SetKey(key string) {
	e.key = key
}

// SetError setter error
func (e *EventContext) SetError(err error) {
	e.err = err
}

// Context get current context
func (e *EventContext) Context() context.Context {
	return e.ctx
}

// WorkerType get worker type
func (e *EventContext) WorkerType() string {
	return e.workerType
}

// Header get header
func (e *EventContext) Header() map[string]interface{} {
	return e.header
}

// Key get key
func (e *EventContext) Key() string {
	return e.key
}

// Topic get topic
func (e *EventContext) Topic() string {
	return e.topic
}

// Message context
func (e *EventContext) Message() []byte {
	return e.buff.Bytes()
}

// Err get error
func (e *EventContext) Err() error {
	return e.err
}

// Read buffer
func (e *EventContext) Read(p []byte) (int, error) {
	return e.buff.Read(p)
}

// Write buffer
func (e *EventContext) Write(p []byte) (int, error) {
	if e.buff == nil {
		e.buff = &bytes.Buffer{}
	}

	return e.buff.Write(p)
}

// WriteString buffer
func (e *EventContext) WriteString(p string) (int, error) {
	if e.buff == nil {
		e.buff = &bytes.Buffer{}
	}

	return e.buff.WriteString(p)
}

// Reset method
func (e *EventContext) Reset() {
	e.buff.Reset()

	e.ctx = nil
	e.workerType = ""
	e.header = nil
	e.topic = ""
	e.key = ""
	e.err = nil
}
