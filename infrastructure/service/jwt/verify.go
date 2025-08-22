package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type VerifyClaims struct {
	Code string
	Id   string
	jwt.RegisteredClaims
}

func NewRegisterClaims(id, code string, exp time.Time) VerifyClaims {
	return VerifyClaims{
		Code: code,
		Id:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   "Đăng ký tài khoản " + id,
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}
