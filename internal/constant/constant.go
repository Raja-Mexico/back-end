package constant

import "errors"

var (
	ErrEmailInvalid      = errors.New("email is invalid")
	ErrEmailAlreadyExist = errors.New("email already exist")
	ErrEmailNotFound     = errors.New("email not found")
	ErrPasswordInvalid   = errors.New("password is invalid, minimum eight characters, at least one letter and one number")
	ErrPasswordIsWrong   = errors.New("password is wrong")
)
