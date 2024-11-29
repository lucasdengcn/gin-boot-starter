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
	//
	Extra map[string]interface{} `json:"extra,omitempty"`
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

func NewProblemValidationDetail(field, detail string, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "ValidationError",
		Title:    "Input validation failed",
		Status:   http.StatusBadRequest,
		Detail:   fmt.Sprintf("Field: %s, error: %s", field, detail),
		Instance: c.Request.RequestURI,
	}
}

func NewProblemBindingDetail(err error, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "BindError",
		Title:    "Failed to bind to struct",
		Status:   http.StatusBadRequest,
		Detail:   err.Error(),
		Instance: c.Request.RequestURI,
	}
}

func NewProblemAuthDetail(err error, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "AuthError",
		Title:    "Failed to authenticate request",
		Status:   http.StatusUnauthorized,
		Detail:   err.Error(),
		Instance: c.Request.RequestURI,
	}
}

func NewProblemACLDetail(err error, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "ACLError",
		Title:    "Failed to authorize request",
		Status:   http.StatusForbidden,
		Detail:   err.Error(),
		Instance: c.Request.RequestURI,
	}
}

func NewUnexpectedDetail(err error, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "InternalError",
		Title:    "Failed to process request",
		Status:   http.StatusInternalServerError,
		Detail:   err.Error(),
		Instance: c.Request.RequestURI,
	}
}

func NewProblemServiceDetail(err *ServiceError, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "ServiceError",
		Title:    "Service failed to process",
		Status:   http.StatusInternalServerError,
		Detail:   err.Error(),
		Instance: c.Request.RequestURI,
		Extra: map[string]interface{}{
			"code": err.Code,
		},
	}
}

func NewProblemRepositoryDetail(err *RepositoryError, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "RepositoryError",
		Title:    "Failed to execute SQL",
		Status:   http.StatusInternalServerError,
		Detail:   err.Error(),
		Instance: c.Request.RequestURI,
		Extra: map[string]interface{}{
			"code": err.Code,
		},
	}
}

func NewProblemSecurityDetail(err *SecurityError, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "SecurityError",
		Title:    "Failed on security check",
		Status:   http.StatusForbidden,
		Detail:   err.Error(),
		Instance: c.Request.RequestURI,
		Extra: map[string]interface{}{
			"code": err.Code,
		},
	}
}

func NewProblem404Detail(err *EntityNotFoundError, c *gin.Context) *ProblemDetails {
	return &ProblemDetails{
		Type:     "NotFoundError",
		Title:    "Record not found",
		Status:   http.StatusNotFound,
		Detail:   err.Error(),
		Instance: c.Request.RequestURI,
		Extra: map[string]interface{}{
			"id": err.ID,
		},
	}
}

////////////

func ResponseAsServiceError(ctx *gin.Context, val any) bool {
	err, ok := val.(*ServiceError)
	if ok {
		ctx.JSON(http.StatusInternalServerError, NewProblemServiceDetail(err, ctx))
		return true
	}
	return false
}

func ResponseAsRepositoryError(ctx *gin.Context, val any) bool {
	err, ok := val.(*RepositoryError)
	if ok {
		ctx.JSON(http.StatusInternalServerError, NewProblemRepositoryDetail(err, ctx))
		return true
	}
	return false
}

func ResponseAsSecurityError(ctx *gin.Context, val any) bool {
	err, ok := val.(*SecurityError)
	if ok {
		ctx.JSON(http.StatusForbidden, NewProblemSecurityDetail(err, ctx))
		return true
	}
	return false
}

func ResponseAs404Error(ctx *gin.Context, val any) bool {
	err, ok := val.(*EntityNotFoundError)
	if ok {
		ctx.JSON(http.StatusNotFound, NewProblem404Detail(err, ctx))
		return true
	}
	return false
}

func ResponseAs500Error(ctx *gin.Context, val any) {
	err, ok := val.(error)
	if ok {
		ctx.JSON(http.StatusInternalServerError, NewUnexpectedDetail(err, ctx))
	} else {
		ctx.JSON(http.StatusInternalServerError, NewUnexpectedDetail(fmt.Errorf("Unexpected error: %v", val), ctx))
	}
}

func ResponseOnError(ctx *gin.Context, val any) {
	if ResponseAsSecurityError(ctx, val) {
		return
	}
	if ResponseAsServiceError(ctx, val) {
		return
	}
	if ResponseAsRepositoryError(ctx, val) {
		return
	}
	if ResponseAs404Error(ctx, val) {
		return
	}
	ResponseAs500Error(ctx, val)
}
