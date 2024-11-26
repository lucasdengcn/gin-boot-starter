package core

import "fmt"

// ServiceError define
type ServiceError struct {
	Code     int
	Message  string
	Instance string
}

// NewServiceError creation
func NewServiceError(code int, message, instance string) *ServiceError {
	return &ServiceError{
		Code:     code,
		Message:  message,
		Instance: instance,
	}
}

// Error returns the error message for the ServiceError type
func (e *ServiceError) Error() string {
	return fmt.Sprintf("Error at %s, %d: %s", e.Instance, e.Code, e.Message)
}

// EntityNotFoundError define
type EntityNotFoundError struct {
	ID       any
	Message  string
	Instance string
}

// NewEntityNotFoundError creation
func NewEntityNotFoundError(id any, message, instance string) *EntityNotFoundError {
	return &EntityNotFoundError{
		ID:       id,
		Message:  message,
		Instance: instance,
	}
}

// Error returns the error message for the EntityNotFoundError type
func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("Error at %s, %s, id: %v", e.Instance, e.ID, e.Message)
}
