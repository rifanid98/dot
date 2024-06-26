// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"
	core "dot/core"

	mock "github.com/stretchr/testify/mock"

	mongodb "dot/infrastructure/v1/persistence/mongodb"

	options "go.mongodb.org/mongo-driver/mongo/options"

	readpref "go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Connect provides a mock function with given fields: ctx
func (_m *Client) Connect(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Database provides a mock function with given fields: name, opts
func (_m *Client) Database(name string, opts ...*options.DatabaseOptions) mongodb.Database {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 mongodb.Database
	if rf, ok := ret.Get(0).(func(string, ...*options.DatabaseOptions) mongodb.Database); ok {
		r0 = rf(name, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongodb.Database)
		}
	}

	return r0
}

// Ping provides a mock function with given fields: ctx, rp
func (_m *Client) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	ret := _m.Called(ctx, rp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *readpref.ReadPref) error); ok {
		r0 = rf(ctx, rp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StartSession provides a mock function with given fields: ic
func (_m *Client) StartSession(ic *core.InternalContext) (mongodb.Session, error) {
	ret := _m.Called(ic)

	var r0 mongodb.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(*core.InternalContext) (mongodb.Session, error)); ok {
		return rf(ic)
	}
	if rf, ok := ret.Get(0).(func(*core.InternalContext) mongodb.Session); ok {
		r0 = rf(ic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongodb.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(*core.InternalContext) error); ok {
		r1 = rf(ic)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
