package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"dot/config"
	"dot/core"
	"dot/core/v1/entity"
	"dot/core/v1/port/auth"
	"dot/pkg/helper"
	"dot/pkg/util"
	"github.com/golang-jwt/jwt"

	portAccount "dot/core/v1/port/account"
	portCache "dot/core/v1/port/cache"
)

var log = util.NewLogger()

type authUsecaseImpl struct {
	accountRepository portAccount.AccountRepository
	cacheRepository   portCache.CacheRepository
	cfg               *config.AppConfig
}

func NewAuthUsecase(
	accountRepository portAccount.AccountRepository,
	cacheRepository portCache.CacheRepository,
	cfg *config.AppConfig,
) auth.AuthUsecase {
	return &authUsecaseImpl{
		accountRepository: accountRepository,
		cacheRepository:   cacheRepository,
		cfg:               cfg,
	}
}

func (uc *authUsecaseImpl) Register(ic *core.InternalContext, account *entity.Account) (*entity.Account, *core.CustomError) {
	acc, cerr := uc.accountRepository.FindAccountByEmail(ic, account.Email)
	if cerr != nil {
		return nil, cerr
	}
	if acc != nil {
		return nil, &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "email already registered",
		}
	}

	cerr = account.SetPassword()
	if cerr != nil {
		return nil, cerr
	}

	account, cerr = uc.accountRepository.InsertAccount(ic, account)
	if cerr != nil {
		return nil, cerr
	}

	// clear password
	account.Password = ""
	return account, cerr
}

func (uc *authUsecaseImpl) Login(ic *core.InternalContext, email, password string) (*entity.Jwt, *core.CustomError) {
	account, cerr := uc.accountRepository.FindAccountByEmail(ic, email)
	if cerr != nil {
		return nil, cerr
	}
	if account == nil {
		return nil, &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "account not found",
		}
	}

	valid := helper.CheckPasswordHash(password, account.Password)
	if !valid {
		return nil, &core.CustomError{
			Code: core.WRONG_PASSWORD,
		}
	}

	claim := &entity.JwtClaim{
		Id:       account.Id,
		Verified: account.Verified,
	}

	return uc.createToken(ic, claim, uc.cfg.JwtSecretKey)
}

func (uc *authUsecaseImpl) RefreshToken(ic *core.InternalContext, accountId string) (*entity.Jwt, *core.CustomError) {
	account, cerr := uc.accountRepository.FindAccountById(ic, accountId)
	if cerr != nil {
		return nil, cerr
	}
	if account == nil {
		return nil, &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "account not found",
		}
	}

	claim := &entity.JwtClaim{
		Id:       accountId,
		Verified: account.Verified,
	}

	return uc.createToken(ic, claim, uc.cfg.JwtSecretKey)
}

func (uc *authUsecaseImpl) RevokeToken(ic *core.InternalContext, accountId string) *core.CustomError {
	redisKey := fmt.Sprintf("active_token::%v", accountId)
	err := uc.cacheRepository.Delete(ic, redisKey)
	if err != nil {
		log.Error(ic.ToContext(), "failed uc.cacheRepository.Delete(ic, redisKey)", err.Error())
		return &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}
	return nil
}

func (uc *authUsecaseImpl) IsActiveToken(ic *core.InternalContext, accountId string, token string) *core.CustomError {
	redisKey := fmt.Sprintf("active_token::%v", accountId)
	data, cerr := uc.cacheRepository.Get(ic, redisKey)
	if cerr != nil {
		log.Error(ic.ToContext(), "failed uc.cacheRepository.Get(ic, redisKey)", cerr.Error())
		return &core.CustomError{
			Code:    core.INVALID_JWT_TOKEN,
			Message: cerr.Error(),
		}
	}

	jwtData := entity.Jwt{}
	if data != "" {
		err := json.Unmarshal([]byte(data), &jwtData)
		if err != nil {
			log.Error(ic.ToContext(), "failed json.Unmarshal([]byte(*data), &jwtData)", err)
			return &core.CustomError{
				Code: core.INVALID_JWT_TOKEN,
			}
		}
	}

	if jwtData.AccessToken == token || jwtData.RefreshToken == token {
		return nil
	}

	return &core.CustomError{
		Code: core.INVALID_JWT_TOKEN,
	}
}

func (uc *authUsecaseImpl) ChangePassword(ic *core.InternalContext, oldPassword string, account *entity.Account) *core.CustomError {
	// Get account info by account id
	acc, cerr := uc.accountRepository.FindAccountById(ic, account.Id)
	if cerr != nil {
		return cerr
	}
	if acc == nil {
		return &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "account not found",
		}
	}

	// Check old password
	cerr = acc.CheckPasword(oldPassword)
	if cerr != nil {
		return cerr
	}

	cerr = acc.SetNewPassword(account.Password)
	if cerr != nil {
		return cerr
	}

	_, cerr = uc.accountRepository.UpdateAccount(ic, acc)
	if cerr != nil {
		return cerr
	}

	return nil
}

func (uc *authUsecaseImpl) createToken(ic *core.InternalContext, claim *entity.JwtClaim, secret string) (*entity.Jwt, *core.CustomError) {
	accessTokenExpiredTime, refreshTokenDuration, refreshTokenExpiredTime := generateTokenTime(uc.cfg)

	mapClaim := jwt.MapClaims{}
	cerr := helper.StringToStruct(helper.DataToString(claim), &mapClaim)
	if cerr != nil {
		log.Error(ic.ToContext(), "failed to convert *entity.JwtClaim to jwt.MapClaims", cerr.Error(), cerr.Errors, cerr.Message)
		return nil, cerr
	}

	token, cerr := generateToken(ic, mapClaim, accessTokenExpiredTime, refreshTokenExpiredTime, secret)
	if cerr != nil {
		return nil, cerr
	}

	cerr = uc.cacheToken(ic, token, claim.Id, &refreshTokenDuration)
	if cerr != nil {
		return nil, cerr
	}

	token.AccessTokenExpired = accessTokenExpiredTime
	token.RefreshTokenExpired = refreshTokenExpiredTime

	return token, nil
}

func (uc *authUsecaseImpl) cacheToken(ic *core.InternalContext, token *entity.Jwt, accountId string, duration *time.Duration) *core.CustomError {
	redisKey := fmt.Sprintf("active_token::%v", accountId)
	redisValue, err := json.Marshal(token)
	if err != nil {
		log.Error(ic.ToContext(), "failed json.Marshal", err.Error())
		return &core.CustomError{
			Code: core.BAD_REQUEST,
		}
	}

	cerr := uc.cacheRepository.Set(ic, redisKey, string(redisValue), duration)
	if cerr != nil {
		log.Error(ic.ToContext(), "failed to set redis data : %v", cerr.Error())
		return cerr
	}

	return nil
}

func generateTokenTime(cfg *config.AppConfig) (int64, time.Duration, int64) {
	accessTokenExpiredTime := time.Now().Add(time.Duration(cfg.JwtAccessTokenExpire) * time.Hour).Unix()
	refreshTokenDuration := time.Duration(cfg.JwtRefreshTokenExpire) * time.Hour
	refreshTokenExpiredTime := time.Now().Add(refreshTokenDuration).Unix()
	return accessTokenExpiredTime, refreshTokenDuration, refreshTokenExpiredTime
}

func generateToken(ic *core.InternalContext, claim jwt.MapClaims, accessTokenExpiredTime int64, refreshTokenExpiredTime int64, secret string) (*entity.Jwt, *core.CustomError) {
	now := time.Now()

	claim["is_refresh"] = false
	claim["iat"] = now.Unix()
	claim["exp"] = accessTokenExpiredTime
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Error(ic.ToContext(), "failed token.SignedString([]byte(secret))", err.Error())
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	claim["is_refresh"] = true
	claim["iat"] = now.Unix()
	claim["exp"] = refreshTokenExpiredTime
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	refreshToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Error(ic.ToContext(), "failed token.SignedString([]byte(config.GetConfig().JwtSecretKey))", err.Error())
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	return &entity.Jwt{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
