// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	retrier "dot/core/v1/port/retrier"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Retrier is an autogenerated mock type for the Retrier type
type Retrier struct {
	mock.Mock
}

// Retry provides a mock function with given fields: effector, retries, delay
func (_m *Retrier) Retry(effector retrier.Effector, retries int, delay time.Duration) retrier.Effector {
	ret := _m.Called(effector, retries, delay)

	var r0 retrier.Effector
	if rf, ok := ret.Get(0).(func(retrier.Effector, int, time.Duration) retrier.Effector); ok {
		r0 = rf(effector, retries, delay)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(retrier.Effector)
		}
	}

	return r0
}

// NewRetrier creates a new instance of Retrier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRetrier(t interface {
	mock.TestingT
	Cleanup(func())
}) *Retrier {
	mock := &Retrier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
