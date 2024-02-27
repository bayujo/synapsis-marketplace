package errors

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFailedClaims       = errors.New("unable to parse claims")
	ErrEmptyCart          = errors.New("cart is empty")
)
