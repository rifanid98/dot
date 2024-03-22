package auth

import (
	"dot/core"
	"dot/core/v1/entity"
)

//go:generate mockery --name AuthUsecase --filename auth_usecase.go --output ./mocks
type AuthUsecase interface {
	Register(ic *core.InternalContext, account *entity.Account) (*entity.Account, *core.CustomError)
	Login(ic *core.InternalContext, email, password string) (*entity.Jwt, *core.CustomError)
	RefreshToken(ic *core.InternalContext, accountId string) (*entity.Jwt, *core.CustomError)
	RevokeToken(ic *core.InternalContext, accountId string) *core.CustomError
	IsActiveToken(ic *core.InternalContext, accountId, token string) *core.CustomError
	ChangePassword(ic *core.InternalContext, oldPassword string, account *entity.Account) *core.CustomError
}
