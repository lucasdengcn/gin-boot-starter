package security

import (
	"gin001/core/logging"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SaveCurrentUser to attache current user to request context
func SaveCurrentUser(c *gin.Context, token *jwt.Token) {
	principle := NewPrinciple(token)
	c.Set(PrincipleContextKey, principle)
	logging.Debug(c).Msgf("Current user is: %s", principle)
}

// CurrentUser return current user attach to request context
func CurrentUser(c *gin.Context) *Principle {
	val, exists := c.Get(PrincipleContextKey)
	if !exists {
		return nil
	}
	return val.(*Principle)
}

// CurrentUser return current user attach to request context
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(PrincipleContextKey)
	return exists
}
