// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	request "github.com/mqdvi-dp/go-common/request"
	mock "github.com/stretchr/testify/mock"
)

// ApiClient is an autogenerated mock type for the ApiClient type
type ApiClient struct {
	mock.Mock
}

// Request provides a mock function with given fields: header, target, url
func (_m *ApiClient) Request(header http.Header, target string, url string) request.MethodInterface {
	ret := _m.Called(header, target, url)

	if len(ret) == 0 {
		panic("no return value specified for Request")
	}

	var r0 request.MethodInterface
	if rf, ok := ret.Get(0).(func(http.Header, string, string) request.MethodInterface); ok {
		r0 = rf(header, target, url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(request.MethodInterface)
		}
	}

	return r0
}

// RequestWithBasicAuth provides a mock function with given fields: header, username, password, target, url
func (_m *ApiClient) RequestWithBasicAuth(header http.Header, username string, password string, target string, url string) request.MethodInterface {
	ret := _m.Called(header, username, password, target, url)

	if len(ret) == 0 {
		panic("no return value specified for RequestWithBasicAuth")
	}

	var r0 request.MethodInterface
	if rf, ok := ret.Get(0).(func(http.Header, string, string, string, string) request.MethodInterface); ok {
		r0 = rf(header, username, password, target, url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(request.MethodInterface)
		}
	}

	return r0
}

// NewApiClient creates a new instance of ApiClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApiClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ApiClient {
	mock := &ApiClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
