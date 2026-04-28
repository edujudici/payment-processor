package domain

import "errors"

var (
	ErrInvalidQuantity = errors.New("invalid quantity must be greater than zero")
	ErrInvalidStatus   = errors.New("invalid status")
)
