package domain

import "errors"

var (
	ErrInvalidAmount = errors.New("invalid amount must be greater than zero")
	ErrInvalidStatus = errors.New("invalid status")
)
