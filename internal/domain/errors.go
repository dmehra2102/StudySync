package domain

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrInavlidInput = errors.New("invalid input")
	ErrConflict     = errors.New("resource conflict")
)
