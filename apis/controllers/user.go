package controllers

import (
	"gin001/apis/models"
	"gin001/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController struct
type UserController struct {
	userService *services.UserService
}

// NewUserController with dependencies services
func NewUserController(UserService *services.UserService) *UserController {
	return &UserController{userService: UserService}
}

// SignUp with user input
func (uc *UserController) SignUp(c *gin.Context) {
	var m models.UserSignUp
	if err := c.ShouldBind(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// call service
	//
	c.JSON(http.StatusCreated, m)
}

// SignIn with user input
func (uc *UserController) SignIn(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
	})
}

// GetUser profile info.
func (uc *UserController) GetUser(c *gin.Context) {
	if uc.userService == nil {
		log.Println("user service is nil")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := uc.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
