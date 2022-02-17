package service

import (
	"errors"
)

var (
	ErrNotSupportGrantType        = errors.New("grant type is not support")
	ErrNotSupportOperation        = errors.New("no support operation")
	ErrInvalidUsernameAndPassword = errors.New("invalid username and password")
	ErrInvalidTokenRequest        = errors.New("invalid token")
	ErrExpiredToken               = errors.New("token is expired")
)


