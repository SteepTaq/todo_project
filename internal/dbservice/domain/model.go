package domain

import (
	"errors"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrTasksNotFound = errors.New("tasks not found")
)
