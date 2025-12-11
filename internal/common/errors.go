package common

import (
	"errors"
)

var (
	ErrNotFound = errors.New("record_not_found")
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Status  int
	Message string
	Errors  []ValidationError
}

func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(errors []ValidationError) *AppError {
	return &AppError{
		Status:  400,
		Message: "Validation failed",
		Errors:  errors,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Status:  404,
		Message: message,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Status:  500,
		Message: message,
	}
}

func NewInvalidRequestBody() *AppError {
	return &AppError{
		Status:  400,
		Message: "Invalid request body",
	}
}

func NewInvalidRequestParams() *AppError {
	return &AppError{
		Status:  400,
		Message: "Invalid request params",
	}
}
