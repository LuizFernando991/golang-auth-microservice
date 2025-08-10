package util

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserExists = errors.New("user already exists")
var ErrRefreshTokenNotFound = errors.New("refresh token not found or expired")
