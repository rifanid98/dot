package auth

import (
	"dot/core/v1/entity"
)

type Login struct {
	Email    string `json:"email" validate:"required,email" example:"dot@email.com"`
	Password string `json:"password" validate:"required,min=6" example:"password"`
}

type Register struct {
	Email    string `json:"email" validate:"required,email" example:"dot@email.com"`
	Password string `json:"password" validate:"required,min=6" example:"password"`
}

func (r *Register) Account() *entity.Account {
	return &entity.Account{
		Email:    r.Email,
		Password: r.Password,
	}
}

type ForgotPassword struct {
	Email    string `json:"email" validate:"required,email" example:"dot@email.com"`
	Otp      string `json:"otp" validate:"required" example:"123456"`
	AppId    string `json:"app_id" validate:"required" swaggerignore:"true"`
	ClientId string `json:"client_id" validate:"required" swaggerignore:"true"`
}

func (r *ForgotPassword) Account() *entity.Account {
	return &entity.Account{
		Email: r.Email,
		Otp:   r.Otp,
	}
}

type ResetPassword struct {
	Email           string `json:"email" validate:",email" example:"dot@email.com"`
	ResetToken      string `json:"reset_token" validate:"required" example:"$2a$14$sEz1WrJXP4pimiJ60WY0Ue6rU1CQjOCmzM7F7G1keIcVPc4tpFnrS"`
	Password        string `json:"password" validate:"required" example:"123456"`
	PasswordConfirm string `json:"password_confirm" validate:"required" example:"123456"`
}

func (r *ResetPassword) Account() *entity.Account {
	return &entity.Account{
		Email:      r.Email,
		ResetToken: r.ResetToken,
		Password:   r.Password,
	}
}

type ChangePassword struct {
	OldPassword     string `json:"old_password" validate:"required" example:"123456"`
	Password        string `json:"password" validate:"required" example:"123456"`
	PasswordConfirm string `json:"password_confirm" validate:"required" example:"123456"`
}

func (r *ChangePassword) Account() *entity.Account {
	return &entity.Account{
		Password: r.Password,
	}
}
