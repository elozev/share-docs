package services

import "errors"

var (
	ErrInvalidId      = errors.New("Invalid ID (should be uuid.v4)")
	ErrFailedToUpdate = errors.New("Failed to update")
)
