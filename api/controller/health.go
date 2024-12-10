package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController interface
type HealthController struct{}

// NewHealthController init
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Status return service status
func (h *HealthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "UP")
}
