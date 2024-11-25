package core

import "fmt"

// ServiceError define
type ServiceError struct {
	Code    int
	Message string
}

// NewServiceError creation
func NewServiceError(code int, message string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
	}
}

// Error returns the error message for the ServiceError type
func (e *ServiceError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// EntityNotFoundError define
type EntityNotFoundError struct {
	ID      any
	Message string
}

// NewEntityNotFoundError creation
func NewEntityNotFoundError(id any, message string) *EntityNotFoundError {
	return &EntityNotFoundError{
		ID:      id,
		Message: message,
	}
}

// Error returns the error message for the EntityNotFoundError type
func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("Error, %s, id: %v", e.ID, e.Message)
}
