package repository

import "errors"

var (
	ErrNotFound      = errors.New("todo not found")
	ErrDatabaseError = errors.New("database error")
)
