package jwt

import (
	"errors"
)

var ErrParseToken = errors.New("parse token error")

var ErrTokenNotFound = errors.New("token not found")
