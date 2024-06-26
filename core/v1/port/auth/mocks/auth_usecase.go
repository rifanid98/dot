// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	core "dot/core"
	entity "dot/core/v1/entity"

	mock "github.com/stretchr/testify/mock"
)

// AuthUsecase is an autogenerated mock type for the AuthUsecase type
type AuthUsecase struct {
	mock.Mock
}

// ChangePassword provides a mock function with given fields: ic, oldPassword, account
func (_m *AuthUsecase) ChangePassword(ic *core.InternalContext, oldPassword string, account *entity.Account) *core.CustomError {
	ret := _m.Called(ic, oldPassword, account)

	var r0 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string, *entity.Account) *core.CustomError); ok {
		r0 = rf(ic, oldPassword, account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.CustomError)
		}
	}

	return r0
}

// IsActiveToken provides a mock function with given fields: ic, accountId, token
func (_m *AuthUsecase) IsActiveToken(ic *core.InternalContext, accountId string, token string) *core.CustomError {
	ret := _m.Called(ic, accountId, token)

	var r0 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string, string) *core.CustomError); ok {
		r0 = rf(ic, accountId, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.CustomError)
		}
	}

	return r0
}

// Login provides a mock function with given fields: ic, email, password
func (_m *AuthUsecase) Login(ic *core.InternalContext, email string, password string) (*entity.Jwt, *core.CustomError) {
	ret := _m.Called(ic, email, password)

	var r0 *entity.Jwt
	var r1 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string, string) (*entity.Jwt, *core.CustomError)); ok {
		return rf(ic, email, password)
	}
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string, string) *entity.Jwt); ok {
		r0 = rf(ic, email, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Jwt)
		}
	}

	if rf, ok := ret.Get(1).(func(*core.InternalContext, string, string) *core.CustomError); ok {
		r1 = rf(ic, email, password)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*core.CustomError)
		}
	}

	return r0, r1
}

// RefreshToken provides a mock function with given fields: ic, accountId
func (_m *AuthUsecase) RefreshToken(ic *core.InternalContext, accountId string) (*entity.Jwt, *core.CustomError) {
	ret := _m.Called(ic, accountId)

	var r0 *entity.Jwt
	var r1 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string) (*entity.Jwt, *core.CustomError)); ok {
		return rf(ic, accountId)
	}
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string) *entity.Jwt); ok {
		r0 = rf(ic, accountId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Jwt)
		}
	}

	if rf, ok := ret.Get(1).(func(*core.InternalContext, string) *core.CustomError); ok {
		r1 = rf(ic, accountId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*core.CustomError)
		}
	}

	return r0, r1
}

// Register provides a mock function with given fields: ic, account
func (_m *AuthUsecase) Register(ic *core.InternalContext, account *entity.Account) (*entity.Account, *core.CustomError) {
	ret := _m.Called(ic, account)

	var r0 *entity.Account
	var r1 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, *entity.Account) (*entity.Account, *core.CustomError)); ok {
		return rf(ic, account)
	}
	if rf, ok := ret.Get(0).(func(*core.InternalContext, *entity.Account) *entity.Account); ok {
		r0 = rf(ic, account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(*core.InternalContext, *entity.Account) *core.CustomError); ok {
		r1 = rf(ic, account)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*core.CustomError)
		}
	}

	return r0, r1
}

// RevokeToken provides a mock function with given fields: ic, accountId
func (_m *AuthUsecase) RevokeToken(ic *core.InternalContext, accountId string) *core.CustomError {
	ret := _m.Called(ic, accountId)

	var r0 *core.CustomError
	if rf, ok := ret.Get(0).(func(*core.InternalContext, string) *core.CustomError); ok {
		r0 = rf(ic, accountId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.CustomError)
		}
	}

	return r0
}

// NewAuthUsecase creates a new instance of AuthUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthUsecase {
	mock := &AuthUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
