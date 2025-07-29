package service

import (
	"errors"
)

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUserExists               = errors.New("user already exists")
	ErrTokenExpired             = errors.New("token has expired")
	ErrInvalidToken             = errors.New("invalid token")
	ErrTokenValidation          = errors.New("token validation error")
	ErrTokenGeneration          = errors.New("failed to generate token")
	ErrUnsupportedSigningMethod = errors.New("unsupported signing method")
	ErrTokenSigning             = errors.New("failed to sign JWT token")

	SqlQuerySyntaxError = errors.New("sql query syntax error")
	NullNotAllowedError = errors.New("sql insert null into not null error")
	SomethingWrongError = errors.New("something wrong")
)
