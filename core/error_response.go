package core

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProblemDetails at https://tools.ietf.org/html/rfc7807
type ProblemDetails struct {
	// Type is a URI reference that identifies the problem type.
	Type string `json:"type"`
	// Title is a short summary of the problem type
	Title string `json:"title"`
	// Status is the HTTP status code for this problem
	Status int `json:"status,omitempty"`
	// Detail is explanation specific to this problem
	Detail string `json:"detail,omitempty"`
	// Instance is a URI reference that identifies this problem.
	Instance string `json:"instance,omitempty"`
}

func NewProblemDetails(statusCode int, problemType, title, detail, instance string) *ProblemDetails {
	// When this member is not present, its value is assumed to be
	// "about:blank".
	if problemType == "" {
		problemType = "about:blank"
	}

	// When "about:blank" is used, the title SHOULD be the same as the
	// recommended HTTP status phrase for that code (e.g., "Not Found" for
	// 404, and so on), although it MAY be localized to suit client
	// preferences (expressed with the Accept-Language request header).
	if problemType == "about:blank" {
		title = http.StatusText(statusCode)
	}

	return &ProblemDetails{
		Type:     problemType,
		Title:    title,
		Status:   statusCode,
		Detail:   detail,
		Instance: instance,
	}
}

// Error returns the error message for the ProblemDetails type
func (e *ProblemDetails) Error() string {
	return fmt.Sprintf("Error: type=%s, title=%s, status=%v, detail=%v, instance=%v", e.Type, e.Title, e.Status, e.Detail, e.Instance)
}

// NewHTTPStatus creates a new ProblemDetails error based just the HTTP Status Code
func NewHTTPStatus(statusCode int) *ProblemDetails {
	return NewProblemDetails(statusCode, "", "", "", "")
}

func NewValidationError(field, detail string, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "ValidationError",
		Title:    "Input validation failed",
		Status:   http.StatusBadRequest,
		Detail:   fmt.Sprintf("Field: %s, error: %s", field, detail),
		Instance: c.Request.RequestURI,
	}
}

func NewBindingError(err error, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "BindError",
		Title:    "Failed to bind to struct",
		Status:   http.StatusBadRequest,
		Detail:   err.Error(),
		Instance: c.Request.RequestURI,
	}
}
