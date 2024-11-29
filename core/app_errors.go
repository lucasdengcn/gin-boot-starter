package core

import "fmt"

// //////////////
// ServiceError define
type ServiceError struct {
	Code     int
	Instance string
	Message  string
}

// NewServiceError creation
func NewServiceError(code int, message, instance string) *ServiceError {
	return &ServiceError{
		Code:     code,
		Instance: instance,
		Message:  message,
	}
}

// Error returns the error message for the ServiceError type
func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s, %s", e.Instance, e.Message)
}

// //////////////
// EntityNotFoundError define
type EntityNotFoundError struct {
	ID       any
	Instance string
	Message  string
}

// NewEntityNotFoundError creation
func NewEntityNotFoundError(id any, message, instance string) *EntityNotFoundError {
	return &EntityNotFoundError{
		ID:       id,
		Instance: instance,
		Message:  message,
	}
}

// Error returns the error message for the EntityNotFoundError type
func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("%s, %s", e.Instance, e.Message)
}

// //////////////
type RepositoryError struct {
	Code     int
	Instance string
	Message  string
}

// NewRepositoryError creation
func NewRepositoryError(code int, message, instance string) *RepositoryError {
	return &RepositoryError{
		Code:     code,
		Instance: instance,
		Message:  message,
	}
}

// Error returns the error message for the ServiceError type
func (e *RepositoryError) Error() string {
	return fmt.Sprintf("%s, %s", e.Instance, e.Message)
}

// //////////////
type SecurityError struct {
	Code     int
	Instance string
	Message  string
}

// NewSecurityError creation
func NewSecurityError(code int, message, instance string) *SecurityError {
	return &SecurityError{
		Code:     code,
		Instance: instance,
		Message:  message,
	}
}

// Error returns the error message for the ServiceError type
func (e *SecurityError) Error() string {
	return fmt.Sprintf("%s, %s", e.Instance, e.Message)
}
