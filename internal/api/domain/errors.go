package domain

import (
	"errors"
)

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyExists  = errors.New("task already exists")
	ErrInvalidInput       = errors.New("invalid input data")
	ErrRequestTimeout     = errors.New("request timeout")
	ErrServiceUnavailable = errors.New("service unavailable")
)
