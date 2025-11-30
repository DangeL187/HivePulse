package domain

import "errors"

var ErrEmailAlreadyExists = errors.New("email already in use")
var ErrUserAlreadyExists = errors.New("email already in use")
var ErrUserNotFound = errors.New("user not found")
