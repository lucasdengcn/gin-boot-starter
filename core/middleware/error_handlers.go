package middleware

import (
	"gin-boot-starter/core/exception"
	"gin-boot-starter/core/logging"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a custom middleware for handling errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				logging.Error(c).Msgf("Recovered from panic: %T, %v", err, err)
				exception.ResponseOnError(c, err)
			}
		}()
		c.Next()
	}
}
