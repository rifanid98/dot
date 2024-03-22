package common

type JwtClaims struct {
	Id       string `json:"id"`
	Verified int    `json:"verified"`
}
