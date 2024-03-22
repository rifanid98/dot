package account

import (
	"dot/core"
	"dot/core/v1/entity"
)

//go:generate mockery --name AccountRepository --filename account_repository.go --output ./mocks
type AccountRepository interface {
	InsertAccount(ic *core.InternalContext, account *entity.Account) (*entity.Account, *core.CustomError)
	FindAccountByEmail(ic *core.InternalContext, email string) (*entity.Account, *core.CustomError)
	FindAccountById(ic *core.InternalContext, accountId string) (*entity.Account, *core.CustomError)
	UpdateAccount(ic *core.InternalContext, account *entity.Account) (*entity.Account, *core.CustomError)
	DeleteAccount(ic *core.InternalContext, id string) *core.CustomError
}
