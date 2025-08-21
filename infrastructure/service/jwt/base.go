package pkgjwt

import (
	serviceJwt "cms-server/domain/service/jwt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtImpl struct {
	secretKey string
}

func NewJWT(secretKey string) serviceJwt.JwtService {
	return &jwtImpl{
		secretKey: secretKey,
	}
}

func (j *jwtImpl) SetSecretKey(secretKey string) serviceJwt.JwtService {
	j.secretKey = secretKey
	return j
}

func (j *jwtImpl) generateToken(data jwt.Claims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, data)
}

func (j *jwtImpl) verifyClaim(token string, data jwt.Claims) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(token, data, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(j.secretKey), nil
	})
	return t, err
}

func (j *jwtImpl) VerifyRegisterToken(token string) (*serviceJwt.VerifyClaims, error) {
	t, err := j.verifyClaim(token, &VerifyClaims{})
	if err != nil {
		return nil, err
	}
	claim, ok := t.Claims.(*VerifyClaims)
	if !ok {
		return nil, ErrParseToken
	}
	return &serviceJwt.VerifyClaims{
		Code: claim.Code,
		Id:   claim.Id,
		RegisteredClaims: serviceJwt.RegisteredClaims{
			ExpiresAt: claim.ExpiresAt.Time,
			Subject:   claim.Subject,
			Audience:  claim.Audience,
			NotBefore: claim.NotBefore.Time,
			IssuedAt:  claim.IssuedAt.Time,
		},
	}, nil
}

func (j *jwtImpl) GenRegisterToken(id, code string, exp time.Time) (string, error) {
	data := NewRegisterClaims(id, code, exp)
	token := j.generateToken(data)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtImpl) GenAuthToken(id, fullName string, exp time.Time) (string, error) {
	data := NewAuthClaims(id, fullName, exp)
	token := j.generateToken(data)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtImpl) VerifyAuthToken(token string) (*serviceJwt.AuthClaims, error) {
	t, err := j.verifyClaim(token, &AuthClaims{})
	if err != nil {
		return nil, err
	}
	claim, ok := t.Claims.(*AuthClaims)
	if !ok {
		return nil, ErrParseToken
	}
	return &serviceJwt.AuthClaims{
		Id:       claim.Id,
		FullName: claim.FullName,
		RegisteredClaims: serviceJwt.RegisteredClaims{
			ExpiresAt: claim.ExpiresAt.Time,
			Subject:   claim.Subject,
			Audience:  claim.Audience,
			NotBefore: claim.NotBefore.Time,
			IssuedAt:  claim.IssuedAt.Time,
		},
	}, nil
}

func (j *jwtImpl) GenForgotPasswordToken(id, fullName string, exp time.Time) (string, error) {
	data := NewForgotClaims(id, fullName, exp)
	token := j.generateToken(data)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtImpl) VerifyForgotPasswordToken(token string) (*serviceJwt.ForgotPasswordClaims, error) {
	t, err := j.verifyClaim(token, &ForgotPasswordClaims{})
	if err != nil {
		return nil, err
	}
	claim, ok := t.Claims.(*ForgotPasswordClaims)
	if !ok {
		return nil, ErrParseToken
	}
	return &serviceJwt.ForgotPasswordClaims{
		Code: claim.Code,
		Id:   claim.Id,
		RegisteredClaims: serviceJwt.RegisteredClaims{
			ExpiresAt: claim.ExpiresAt.Time,
			Subject:   claim.Subject,
			Audience:  claim.Audience,
			NotBefore: claim.NotBefore.Time,
			IssuedAt:  claim.IssuedAt.Time,
		},
	}, nil
}
