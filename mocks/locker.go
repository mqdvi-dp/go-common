// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Locker is an autogenerated mock type for the Locker type
type Locker struct {
	mock.Mock
}

// HasBeenLocked provides a mock function with given fields: key
func (_m *Locker) HasBeenLocked(key string) bool {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for HasBeenLocked")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsLocked provides a mock function with given fields: key
func (_m *Locker) IsLocked(key string) bool {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for IsLocked")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Reset provides a mock function with given fields: key
func (_m *Locker) Reset(key string) {
	_m.Called(key)
}

// Unlock provides a mock function with given fields: key
func (_m *Locker) Unlock(key string) {
	_m.Called(key)
}

// NewLocker creates a new instance of Locker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLocker(t interface {
	mock.TestingT
	Cleanup(func())
}) *Locker {
	mock := &Locker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
