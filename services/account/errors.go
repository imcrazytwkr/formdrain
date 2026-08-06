package account

import "errors"

var ErrEmptyPassword = errors.New("password is empty")
var ErrPasswordTooLong = errors.New("password is too long")
var ErrInvalidHash = errors.New("invalid password hash")
var ErrInvalidCredentials = errors.New("invalid email or password")
