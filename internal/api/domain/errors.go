package domain

import (
	"errors"
)

// Ошибки предметной области
var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyExists  = errors.New("task already exists")
	ErrInvalidInput       = errors.New("invalid input data")
	ErrRequestTimeout     = errors.New("request timeout")
	ErrServiceUnavailable = errors.New("service unavailable")
)
