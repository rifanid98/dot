package entity

import (
	"github.com/golang-jwt/jwt"
)

type JwtClaim struct {
	// General claim
	Id        string `json:"id"`
	Verified  int    `json:"verified"`
	IsRefresh bool   `json:"is_refresh"`
	jwt.StandardClaims
}

type Jwt struct {
	AccessToken         string `json:"access_token"`
	AccessTokenExpired  int64  `json:"access_token_expired"`
	RefreshToken        string `json:"refresh_token"`
	RefreshTokenExpired int64  `json:"refresh_token_expired"`
}
