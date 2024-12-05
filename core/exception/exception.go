package exception

// to avoid conflict go errors package

import (
	"fmt"
	"gin-boot-starter/core/correlation"

	"github.com/gin-gonic/gin"
)

// ServiceError define
type ServiceError struct {
	Errno         string
	CorrelationId string
	Message       string
}

// NewServiceError creation
func NewServiceError(ctx *gin.Context, errno string, message string) *ServiceError {
	return &ServiceError{
		Errno:         errno,
		CorrelationId: correlation.CorrelationId(ctx),
		Message:       message,
	}
}

// Error returns the error message for the ServiceError type
func (e *ServiceError) Error() string {
	return fmt.Sprintf("errno: %s, correlationId: %s, %s", e.Errno, e.CorrelationId, e.Message)
}

func (s *ServiceError) Is(err error) bool {
	_, ok := err.(*ServiceError)
	if !ok {
		return false
	}
	return true
}

// EntityNotFoundError define
type EntityNotFoundError struct {
	ID            any
	CorrelationId string
	Message       string
}

// NewEntityNotFoundError creation
func NewEntityNotFoundError(ctx *gin.Context, id any, message string) *EntityNotFoundError {
	return &EntityNotFoundError{
		ID:            id,
		CorrelationId: correlation.CorrelationId(ctx),
		Message:       message,
	}
}

// Error returns the error message for the EntityNotFoundError type
func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("Id: %d, correlationId: %s, %s", e.ID, e.CorrelationId, e.Message)
}

func (s *EntityNotFoundError) Is(err error) bool {
	_, ok := err.(*EntityNotFoundError)
	if !ok {
		return false
	}
	return true
}

// RepositoryError
type RepositoryError struct {
	Errno         string
	CorrelationId string
	Message       string
}

// NewRepositoryError creation
func NewRepositoryError(ctx *gin.Context, errno string, message string) *RepositoryError {
	return &RepositoryError{
		Errno:         errno,
		CorrelationId: correlation.CorrelationId(ctx),
		Message:       message,
	}
}

// Error returns the error message for the ServiceError type
func (e *RepositoryError) Error() string {
	return fmt.Sprintf("errno: %s, correlationId: %s, %s", e.Errno, e.CorrelationId, e.Message)
}

func (s *RepositoryError) Is(err error) bool {
	_, ok := err.(*RepositoryError)
	if !ok {
		return false
	}
	return true
}

// ACLError
type ACLError struct {
	Errno         string
	CorrelationId string
	Message       string
}

// NewACLError creation
func NewACLError(ctx *gin.Context, errno string, message string) *ACLError {
	return &ACLError{
		Errno:         errno,
		CorrelationId: correlation.CorrelationId(ctx),
		Message:       message,
	}
}

// Error returns the error message for the ServiceError type
func (e *ACLError) Error() string {
	return fmt.Sprintf("errno: %s, correlationId: %s, %s", e.Errno, e.CorrelationId, e.Message)
}

func (s *ACLError) Is(err error) bool {
	_, ok := err.(*ACLError)
	if !ok {
		return false
	}
	return true
}

// AuthError
type AuthError struct {
	Errno         string
	CorrelationId string
	Message       string
}

// NewAuthError creation
func NewAuthError(ctx *gin.Context, errno string, message string) *AuthError {
	return &AuthError{
		Errno:         errno,
		CorrelationId: correlation.CorrelationId(ctx),
		Message:       message,
	}
}

// Error returns the error message for the ServiceError type
func (e *AuthError) Error() string {
	return fmt.Sprintf("errno: %s, correlationId: %s, %s", e.Errno, e.CorrelationId, e.Message)
}

func (s *AuthError) Is(err error) bool {
	_, ok := err.(*AuthError)
	if !ok {
		return false
	}
	return true
}

// ValidationError
type ValidationError struct {
	Errno         string
	CorrelationId string
	Message       string
}

// NewValidationError creation
func NewValidationError(ctx *gin.Context, errno string, message string) *ValidationError {
	return &ValidationError{
		Errno:         errno,
		CorrelationId: correlation.CorrelationId(ctx),
		Message:       message,
	}
}

// Error returns the error message for the ServiceError type
func (e *ValidationError) Error() string {
	return fmt.Sprintf("errno: %s, correlationId: %s, %s", e.Errno, e.CorrelationId, e.Message)
}

func (s *ValidationError) Is(err error) bool {
	_, ok := err.(*ValidationError)
	if !ok {
		return false
	}
	return true
}
