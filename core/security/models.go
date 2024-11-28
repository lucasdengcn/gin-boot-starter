package security

import (
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

var PrincipleContextKey = "_principle_context_key_"

// Principle to indicate current session user
type Principle struct {
	Token  *jwt.Token
	ID     string
	Claims *AuthClaims
}

func NewPrinciple(token *jwt.Token) *Principle {
	//
	claims := token.Claims.(*AuthClaims)
	id, _ := claims.GetSubject()
	//
	return &Principle{
		Token:  token,
		ID:     id,
		Claims: claims,
	}
}

func (p *Principle) GetID() uint {
	val, _ := strconv.Atoi(p.ID)
	return uint(val)
}

func (p *Principle) GetRoles() []string {
	if p.Claims == nil {
		return nil
	}
	return p.Claims.Roles
}

func (p *Principle) GetGroups() []string {
	if p.Claims == nil {
		return nil
	}
	return p.Claims.Groups
}

func (p *Principle) String() string {
	return fmt.Sprintf("sub:%s, roles:%s, groups:%s", p.ID, p.GetRoles(), p.GetGroups())
}

// AuthClaims for authentication and authorization
type AuthClaims struct {
	Roles  []string `json:"roles"`
	Groups []string `json:"groups"`
	jwt.RegisteredClaims
}
