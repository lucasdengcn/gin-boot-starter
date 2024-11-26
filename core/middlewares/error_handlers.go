package middlewares

import (
	"errors"
	"fmt"
	"gin001/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ErrorHandler is a custom middleware for handling errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				fmt.Println("Recovered from panic:", err)
				url := c.Request.RequestURI
				realError, ok := err.(error)
				//
				if !ok {
					log.Error().Msgf("Unexpected Error: %v, %v", url, err)
					c.JSON(http.StatusInternalServerError, core.NewProblemDetails(500, "Unexpected", "Unexpected Error", fmt.Sprintf("%v", err), url))
					return
				}
				log.Error().Err(realError).Msgf("Request Url: %v", url)
				if errors.Is(realError, &core.ServiceError{}) {
					// Return an error response to the client
					c.JSON(http.StatusInternalServerError, core.NewProblemDetails(500, "ServiceError", "Service Error", realError.Error(), url))
				} else if errors.Is(realError, &core.EntityNotFoundError{}) {
					// Return an error response to the client
					c.JSON(http.StatusNotFound, core.NewProblemDetails(404, "NotFound", "Resource NotFound Error", realError.Error(), url))
				} else {
					c.JSON(http.StatusInternalServerError, core.NewProblemDetails(500, "Unexpected", "Unexpected Error", realError.Error(), url))
				}
			}
		}()
		c.Next()
	}
}
