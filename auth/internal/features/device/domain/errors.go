package domain

import "errors"

var ErrDeviceAlreadyExists = errors.New("device ID already in use")
var ErrDeviceNotFound = errors.New("device not found")
