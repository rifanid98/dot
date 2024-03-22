package auth

import (
	"dot/pkg/helper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"

	"dot/config"
	"dot/core"
	"dot/core/v1/entity"

	accountMocks "dot/core/v1/port/account/mocks"
	cacheMocks "dot/core/v1/port/cache/mocks"
	commonMocks "dot/core/v1/port/common/mocks"
)

type mocks struct {
	accountRepository accountMocks.AccountRepositoryMock
	cacheRepository   cacheMocks.CacheRepositoryMock
}

func Test_authUsecaseImpl_Register(t *testing.T) {
	type fields struct {
		accountRepository *accountMocks.AccountRepository
		cacheRepository   *cacheMocks.CacheRepository
		cfg               *config.AppConfig
	}
	type args struct {
		ic      *core.InternalContext
		account *entity.Account
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mocks  mocks
		want   *entity.Account
		want1  *core.CustomError
	}{
		{
			name: "register account must be success",
			fields: fields{
				accountRepository: new(accountMocks.AccountRepository),
				cacheRepository:   new(cacheMocks.CacheRepository),
				cfg:               &config.AppConfig{},
			},
			args: args{
				ic:      core.NewInternalContext(uuid.NewString()),
				account: &entity.Account{},
			},
			mocks: mocks{
				accountRepository: accountMocks.AccountRepositoryMock{
					FindAccountByEmail:    nil,
					FindAccountByEmailErr: nil,
					InsertAccount:         &entity.Account{},
					InsertAccountErr:      nil,
				},
			},
			want:  &entity.Account{},
			want1: nil,
		},
		{
			name: "register account, but email already registered",
			fields: fields{
				accountRepository: new(accountMocks.AccountRepository),
				cacheRepository:   new(cacheMocks.CacheRepository),
				cfg:               &config.AppConfig{},
			},
			args: args{
				ic:      core.NewInternalContext(uuid.NewString()),
				account: &entity.Account{},
			},
			mocks: mocks{
				accountRepository: accountMocks.AccountRepositoryMock{
					FindAccountByEmail:    &entity.Account{},
					FindAccountByEmailErr: nil,
					InsertAccount:         nil,
					InsertAccountErr:      nil,
				},
			},
			want: nil,
			want1: &core.CustomError{
				Code:    core.UNPROCESSABLE_ENTITY,
				Message: "email already registered",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepository := tt.fields.accountRepository
			cacheRepository := tt.fields.cacheRepository
			cfg := tt.fields.cfg
			mocks := tt.mocks

			uc := &authUsecaseImpl{
				accountRepository: accountRepository,
				cacheRepository:   cacheRepository,
				cfg:               cfg,
			}

			accountRepository.On("FindAccountByEmail", mock.Anything, mock.Anything).Return(mocks.accountRepository.FindAccountByEmail, mocks.accountRepository.FindAccountByEmailErr)
			accountRepository.On("InsertAccount", mock.Anything, mock.Anything).Return(mocks.accountRepository.InsertAccount, mocks.accountRepository.InsertAccountErr)

			got, got1 := uc.Register(tt.args.ic, tt.args.account)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Register() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_authUsecaseImpl_Login(t *testing.T) {
	password, _ := helper.HashPassword("password")
	type fields struct {
		accountRepository *accountMocks.AccountRepository
		cacheRepository   *cacheMocks.CacheRepository
		transaction       *commonMocks.Transaction
		cfg               *config.AppConfig
	}
	type args struct {
		ic       *core.InternalContext
		email    string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mocks  mocks
		want   *entity.Jwt
		want1  *core.CustomError
	}{
		{
			name: "login must be success",
			fields: fields{
				accountRepository: new(accountMocks.AccountRepository),
				cacheRepository:   new(cacheMocks.CacheRepository),
				cfg:               &config.AppConfig{},
			},
			args: args{
				ic:       core.NewInternalContext(uuid.NewString()),
				email:    "email",
				password: "password",
			},
			mocks: mocks{
				accountRepository: accountMocks.AccountRepositoryMock{
					FindAccountByEmail: &entity.Account{
						Password: password,
					},
					FindAccountByEmailErr: nil,
				},
				cacheRepository: cacheMocks.CacheRepositoryMock{
					Set: nil,
				},
			},
			want:  nil,
			want1: nil,
		},
		{
			name: "login failed, invalid password",
			fields: fields{
				accountRepository: new(accountMocks.AccountRepository),
				cacheRepository:   new(cacheMocks.CacheRepository),
				cfg:               &config.AppConfig{},
			},
			args: args{
				ic:       core.NewInternalContext(uuid.NewString()),
				email:    "email",
				password: "password-new",
			},
			mocks: mocks{
				accountRepository: accountMocks.AccountRepositoryMock{
					FindAccountByEmail: &entity.Account{
						Password: password,
					},
					FindAccountByEmailErr: nil,
				},
				cacheRepository: cacheMocks.CacheRepositoryMock{
					Set: nil,
				},
			},
			want: nil,
			want1: &core.CustomError{
				Code: core.WRONG_PASSWORD,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepository := tt.fields.accountRepository
			cacheRepository := tt.fields.cacheRepository
			cfg := tt.fields.cfg
			mocks := tt.mocks

			uc := &authUsecaseImpl{
				accountRepository: accountRepository,
				cacheRepository:   cacheRepository,
				cfg:               cfg,
			}

			accountRepository.On("FindAccountByEmail", mock.Anything, mock.Anything).Return(mocks.accountRepository.FindAccountByEmail, mocks.accountRepository.FindAccountByEmailErr)
			accountRepository.On("InsertAccount", mock.Anything, mock.Anything).Return(mocks.accountRepository.InsertAccount, mocks.accountRepository.InsertAccountErr)
			cacheRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mocks.cacheRepository.Set)

			got, got1 := uc.Login(tt.args.ic, tt.args.email, tt.args.password)
			if got != nil {
				if got.AccessToken == "" && got.AccessTokenExpired < 1 && got.RefreshToken == "" && got.RefreshTokenExpired < 1 {
					t.Errorf("Login() got = %v, want got.AccessToken, got.AccessTokenExpired, got.RefreshToken & got.RefreshTokenExpired is not empty or zero", got)
				}
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Login() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_authUsecaseImpl_RefreshToken(t *testing.T) {
	type fields struct {
		accountRepository *accountMocks.AccountRepository
		cacheRepository   *cacheMocks.CacheRepository
		transaction       *commonMocks.Transaction
		cfg               *config.AppConfig
	}
	type args struct {
		ic        *core.InternalContext
		accountId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mocks  mocks
		want   *entity.Jwt
		want1  *core.CustomError
	}{
		{
			name: "account relogin must be success",
			fields: fields{
				accountRepository: new(accountMocks.AccountRepository),
				cacheRepository:   new(cacheMocks.CacheRepository),
				cfg:               &config.AppConfig{},
			},
			args: args{
				ic:        core.NewInternalContext(uuid.NewString()),
				accountId: uuid.NewString(),
			},
			mocks: mocks{
				accountRepository: accountMocks.AccountRepositoryMock{
					FindAccountById:    &entity.Account{},
					FindAccountByIdErr: nil,
				},
				cacheRepository: cacheMocks.CacheRepositoryMock{
					Set: nil,
				},
			},
			want:  nil,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepository := tt.fields.accountRepository
			cacheRepository := tt.fields.cacheRepository
			cfg := tt.fields.cfg
			mocks := tt.mocks

			uc := &authUsecaseImpl{
				accountRepository: accountRepository,
				cacheRepository:   cacheRepository,
				cfg:               cfg,
			}

			accountRepository.On("FindAccountById", mock.Anything, mock.Anything).Return(mocks.accountRepository.FindAccountById, mocks.accountRepository.FindAccountByIdErr)
			cacheRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mocks.cacheRepository.Set)

			got, got1 := uc.RefreshToken(tt.args.ic, tt.args.accountId)
			if got != nil {
				if got.AccessToken == "" && got.AccessTokenExpired < 1 && got.RefreshToken == "" && got.RefreshTokenExpired < 1 {
					t.Errorf("RefreshToken() got = %v, want got.AccessToken, got.AccessTokenExpired, got.RefreshToken & got.RefreshTokenExpired is not empty or zero", got)
				}
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RefreshToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_authUsecaseImpl_RevokeToken(t *testing.T) {
	type fields struct {
		accountRepository *accountMocks.AccountRepository
		cacheRepository   *cacheMocks.CacheRepository
		transaction       *commonMocks.Transaction
		cfg               *config.AppConfig
	}
	type args struct {
		ic        *core.InternalContext
		accountId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mocks  mocks
		want   *core.CustomError
	}{
		{
			name: "account logout must be success",
			fields: fields{
				accountRepository: new(accountMocks.AccountRepository),
				cacheRepository:   new(cacheMocks.CacheRepository),
				cfg:               &config.AppConfig{},
			},
			args: args{
				ic:        core.NewInternalContext(uuid.NewString()),
				accountId: uuid.NewString(),
			},
			mocks: mocks{
				cacheRepository: cacheMocks.CacheRepositoryMock{
					Delete: nil,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepository := tt.fields.accountRepository
			cacheRepository := tt.fields.cacheRepository
			cfg := tt.fields.cfg
			mocks := tt.mocks

			uc := &authUsecaseImpl{
				accountRepository: accountRepository,
				cacheRepository:   cacheRepository,
				cfg:               cfg,
			}

			cacheRepository.On("Delete", mock.Anything, mock.Anything).Return(mocks.cacheRepository.Delete)

			if got := uc.RevokeToken(tt.args.ic, tt.args.accountId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RevokeToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authUsecaseImpl_IsActiveToken(t *testing.T) {
	type fields struct {
		accountRepository *accountMocks.AccountRepository
		cacheRepository   *cacheMocks.CacheRepository
		transaction       *commonMocks.Transaction
		cfg               *config.AppConfig
	}
	type args struct {
		ic        *core.InternalContext
		accountId string
		token     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mocks  mocks
		want   *core.CustomError
	}{
		{
			name: "check token active status must be success",
			fields: fields{
				accountRepository: new(accountMocks.AccountRepository),
				cacheRepository:   new(cacheMocks.CacheRepository),
				cfg:               &config.AppConfig{},
			},
			args: args{
				ic:        core.NewInternalContext(uuid.NewString()),
				accountId: uuid.NewString(),
				token:     "token",
			},
			mocks: mocks{
				cacheRepository: cacheMocks.CacheRepositoryMock{
					Get: helper.DataToString(entity.Jwt{
						AccessToken:         "token",
						AccessTokenExpired:  0,
						RefreshToken:        "token",
						RefreshTokenExpired: 0,
					}),
					GetErr: nil,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepository := tt.fields.accountRepository
			cacheRepository := tt.fields.cacheRepository
			cfg := tt.fields.cfg
			mocks := tt.mocks

			uc := &authUsecaseImpl{
				accountRepository: accountRepository,
				cacheRepository:   cacheRepository,
				cfg:               cfg,
			}

			cacheRepository.On("Get", mock.Anything, mock.Anything).Return(mocks.cacheRepository.Get, mocks.cacheRepository.GetErr)

			if got := uc.IsActiveToken(tt.args.ic, tt.args.accountId, tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsActiveToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authUsecaseImpl_ChangePassword(t *testing.T) {
	passwordOld, _ := helper.HashPassword("password-old")
	type fields struct {
		accountRepository *accountMocks.AccountRepository
		cacheRepository   *cacheMocks.CacheRepository
		transaction       *commonMocks.Transaction
		cfg               *config.AppConfig
	}
	type args struct {
		ic          *core.InternalContext
		oldPassword string
		account     *entity.Account
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mocks  mocks
		want   *core.CustomError
	}{
		{
			name: "change password must be success",
			fields: fields{
				accountRepository: new(accountMocks.AccountRepository),
				cacheRepository:   new(cacheMocks.CacheRepository),
				cfg:               &config.AppConfig{},
			},
			args: args{
				ic:          core.NewInternalContext(uuid.NewString()),
				oldPassword: "password-old",
				account: &entity.Account{
					Password: "password-new",
				},
			},
			mocks: mocks{
				accountRepository: accountMocks.AccountRepositoryMock{
					FindAccountById: &entity.Account{
						Password: passwordOld,
					},
					FindAccountByIdErr: nil,
				},
			},
			want: nil,
		},
		{
			name: "change password, but old password was incorrect",
			fields: fields{
				accountRepository: new(accountMocks.AccountRepository),
				cacheRepository:   new(cacheMocks.CacheRepository),
				cfg:               &config.AppConfig{},
			},
			args: args{
				ic:          core.NewInternalContext(uuid.NewString()),
				oldPassword: "password-incorrect",
				account: &entity.Account{
					Password: "password-new",
				},
			},
			mocks: mocks{
				accountRepository: accountMocks.AccountRepositoryMock{
					FindAccountById: &entity.Account{
						Password: passwordOld,
					},
					FindAccountByIdErr: nil,
				},
			},
			want: &core.CustomError{
				Code:    core.WRONG_PASSWORD,
				Message: "wrong password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepository := tt.fields.accountRepository
			cacheRepository := tt.fields.cacheRepository
			cfg := tt.fields.cfg
			mocks := tt.mocks

			uc := &authUsecaseImpl{
				accountRepository: accountRepository,
				cacheRepository:   cacheRepository,
				cfg:               cfg,
			}

			accountRepository.On("FindAccountById", mock.Anything, mock.Anything).Return(mocks.accountRepository.FindAccountById, mocks.accountRepository.FindAccountByIdErr)
			accountRepository.On("UpdateAccount", mock.Anything, mock.Anything).Return(mocks.accountRepository.UpdateAccount, mocks.accountRepository.UpdateAccountErr)

			if got := uc.ChangePassword(tt.args.ic, tt.args.oldPassword, tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
