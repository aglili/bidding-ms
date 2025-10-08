package repository

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrDatabaseError = errors.New("database error")
)
