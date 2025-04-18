// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Tracer is an autogenerated mock type for the Tracer type
type Tracer struct {
	mock.Mock
}

// Context provides a mock function with given fields:
func (_m *Tracer) Context() context.Context {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Context")
	}

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// Debug provides a mock function with given fields: key, args
func (_m *Tracer) Debug(key string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, key)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Finish provides a mock function with given fields:
func (_m *Tracer) Finish() {
	_m.Called()
}

// Log provides a mock function with given fields: key, args
func (_m *Tracer) Log(key string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, key)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// NewContext provides a mock function with given fields:
func (_m *Tracer) NewContext() context.Context {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewContext")
	}

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// SetError provides a mock function with given fields: err
func (_m *Tracer) SetError(err error) {
	_m.Called(err)
}

// SetTag provides a mock function with given fields: key, value
func (_m *Tracer) SetTag(key string, value interface{}) {
	_m.Called(key, value)
}

// Tags provides a mock function with given fields:
func (_m *Tracer) Tags() map[string]interface{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Tags")
	}

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// NewTracer creates a new instance of Tracer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTracer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Tracer {
	mock := &Tracer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
