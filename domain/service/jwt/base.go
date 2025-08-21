package serviceJwt

import "time"

type JwtService interface {
	GenRegisterToken(id, code string, exp time.Time) (string, error)
	VerifyRegisterToken(token string) (*VerifyClaims, error)
	GenAuthToken(id, fullName string, exp time.Time) (string, error)
	VerifyAuthToken(token string) (*AuthClaims, error)
	GenForgotPasswordToken(id, code string, exp time.Time) (string, error)
	VerifyForgotPasswordToken(token string) (*ForgotPasswordClaims, error)
}

type RegisteredClaims struct {
	ExpiresAt time.Time
	Subject   string
	Audience  []string
	NotBefore time.Time
	IssuedAt  time.Time
}

type VerifyClaims struct {
	Code string
	Id   string
	RegisteredClaims
}

type ForgotPasswordClaims struct {
	Code string
	Id   string
	RegisteredClaims
}

type AuthClaims struct {
	Id       string
	FullName string
	RegisteredClaims
}
