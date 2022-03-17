package auth

import "errors"

var (
	ErrorUserAlreadyExist = errors.New("user already exist")
	ErrorWrongAuthData    = errors.New("incorrect login or password")
)
