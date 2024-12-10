package middleware

import (
	"fmt"
	"gin-boot-starter/core/exception"
	"gin-boot-starter/core/logging"
	"gin-boot-starter/core/security"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	bearer = "Bearer "
)

func AuthJwtHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			tokenString, err := c.Cookie("X-Authorization")
			if tokenString == "" || err != nil {
				c.JSON(http.StatusUnauthorized, exception.NewProblemAuthDetail(fmt.Errorf("Authorization Required"), c))
				c.Abort()
				return
			}
		}
		if !strings.HasPrefix(tokenString, bearer) {
			c.JSON(http.StatusUnauthorized, exception.NewProblemAuthDetail(fmt.Errorf("Authorization Header Invalid"), c))
			c.Abort()
			return
		}
		//
		tokenString = strings.TrimPrefix(tokenString, bearer)
		//
		token, err := jwt.ParseWithClaims(tokenString, &security.AuthClaims{}, security.PublicJwtKeyfuncCtx(c))
		if err != nil || !token.Valid {
			logging.Error(c).Err(err).Msgf("Token Invalid:%s", tokenString)
			c.JSON(http.StatusUnauthorized, exception.NewProblemAuthDetail(fmt.Errorf("Token Invalid"), c))
			c.Abort()
			return
		}
		//
		security.SaveCurrentUser(c, token)
		// call next handler
		c.Next()
	}
}
