package mocks

import (
	"dot/core"
	"dot/core/v1/entity"
)

type AccountRepositoryMock struct {
	FindAccountByAccountId    *entity.Account
	FindAccountByAccountIdErr *core.CustomError
	FindAccountByEmail        *entity.Account
	FindAccountByEmailErr     *core.CustomError
	FindAccountById           *entity.Account
	FindAccountByIdErr        *core.CustomError
	FindAccountByUsername     *entity.Account
	FindAccountByUsernameErr  *core.CustomError
	FindAccountsActivation    []entity.Account
	FindAccountsActivationErr *core.CustomError
	GetAccountsExclude        []entity.Account
	GetAccountsExcludeTotal   int32
	GetAccountsExcludeErr     *core.CustomError
	InsertAccount             *entity.Account
	InsertAccountErr          *core.CustomError
	InsertAccounts            []entity.Account
	InsertAccountsErr         *core.CustomError
	SetAccountOtp             *core.CustomError
	UpdateAccount             *entity.Account
	UpdateAccountErr          *core.CustomError
}

type AccountUsecaseMock struct {
	AccountAction          *core.CustomError
	AccountActivate        map[string]interface{}
	AccountActivateErr     *core.CustomError
	AccountActivationCheck *core.CustomError
	AccountGet             *entity.Account
	AccountGetErr          *core.CustomError
	AccountList            []entity.Account
	AccountListTotal       int32
	AccountListErr         *core.CustomError
}
