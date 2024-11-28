package middlewares

import (
	"fmt"
	"gin001/core"
	"gin001/core/logging"
	"gin001/core/security"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	bearer = "Bearer "
)

func AuthJwtHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, core.NewProblemAuthDetail(fmt.Errorf("Authorization Required"), c))
			c.Abort()
			return
		}
		if !strings.HasPrefix(tokenString, bearer) {
			c.JSON(http.StatusUnauthorized, core.NewProblemAuthDetail(fmt.Errorf("Authorization Header Invalid"), c))
			c.Abort()
			return
		}
		//
		tokenString = strings.TrimPrefix(tokenString, bearer)
		//
		token, err := jwt.ParseWithClaims(tokenString, &security.AuthClaims{}, security.PublicJwtKeyfuncCtx(c))
		if err != nil || !token.Valid {
			logging.Error(c).Err(err).Msgf("Token Invalid:%s", tokenString)
			c.JSON(http.StatusUnauthorized, core.NewProblemAuthDetail(fmt.Errorf("Token Invalid"), c))
			c.Abort()
			return
		}
		//
		security.SaveCurrentUser(c, token)
		// call next handler
		c.Next()
	}
}
