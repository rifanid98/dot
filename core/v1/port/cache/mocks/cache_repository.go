// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	core "dot/core"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// CacheRepository is an autogenerated mock type for the CacheRepository type
type CacheRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ic, key
func (_m *CacheRepository) Delete(ic *core.InternalContext, key string) *core.CustomError {
	ret := _m.Called(ic, key)

	var r0 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string) *core.CustomError); ok {
		r0 = rf(ic, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.CustomError)
		}
	}

	return r0
}

// Get provides a mock function with given fields: ic, key
func (_m *CacheRepository) Get(ic *core.InternalContext, key string) (string, *core.CustomError) {
	ret := _m.Called(ic, key)

	var r0 string
	var r1 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string) (string, *core.CustomError)); ok {
		return rf(ic, key)
	}
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string) string); ok {
		r0 = rf(ic, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*core.InternalContext, string) *core.CustomError); ok {
		r1 = rf(ic, key)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*core.CustomError)
		}
	}

	return r0, r1
}

// HSet provides a mock function with given fields: ic, key, value, expiration
func (_m *CacheRepository) HSet(ic *core.InternalContext, key string, value map[string]interface{}, expiration time.Duration) *core.CustomError {
	ret := _m.Called(ic, key, value, expiration)

	var r0 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string, map[string]interface{}, time.Duration) *core.CustomError); ok {
		r0 = rf(ic, key, value, expiration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.CustomError)
		}
	}

	return r0
}

// Publish provides a mock function with given fields: ic, channel, data
func (_m *CacheRepository) Publish(ic *core.InternalContext, channel string, data string) *core.CustomError {
	ret := _m.Called(ic, channel, data)

	var r0 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string, string) *core.CustomError); ok {
		r0 = rf(ic, channel, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.CustomError)
		}
	}

	return r0
}

// Set provides a mock function with given fields: ic, key, value, expiration
func (_m *CacheRepository) Set(ic *core.InternalContext, key string, value string, expiration *time.Duration) *core.CustomError {
	ret := _m.Called(ic, key, value, expiration)

	var r0 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string, string, *time.Duration) *core.CustomError); ok {
		r0 = rf(ic, key, value, expiration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.CustomError)
		}
	}

	return r0
}

// NewCacheRepository creates a new instance of CacheRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCacheRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CacheRepository {
	mock := &CacheRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
