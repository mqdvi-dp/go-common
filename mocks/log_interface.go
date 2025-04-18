// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	logrus "github.com/sirupsen/logrus"
	mock "github.com/stretchr/testify/mock"
)

// logInterface is an autogenerated mock type for the logInterface type
type logInterface struct {
	mock.Mock
}

// Debug provides a mock function with given fields: args
func (_m *logInterface) Debug(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Debugf provides a mock function with given fields: format, args
func (_m *logInterface) Debugf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Error provides a mock function with given fields: args
func (_m *logInterface) Error(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Errorf provides a mock function with given fields: format, args
func (_m *logInterface) Errorf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Fatal provides a mock function with given fields: args
func (_m *logInterface) Fatal(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Fatalf provides a mock function with given fields: format, args
func (_m *logInterface) Fatalf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Println provides a mock function with given fields: args
func (_m *logInterface) Println(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// WithField provides a mock function with given fields: key, value
func (_m *logInterface) WithField(key string, value interface{}) *logrus.Entry {
	ret := _m.Called(key, value)

	if len(ret) == 0 {
		panic("no return value specified for WithField")
	}

	var r0 *logrus.Entry
	if rf, ok := ret.Get(0).(func(string, interface{}) *logrus.Entry); ok {
		r0 = rf(key, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*logrus.Entry)
		}
	}

	return r0
}

// newLogInterface creates a new instance of logInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newLogInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *logInterface {
	mock := &logInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
