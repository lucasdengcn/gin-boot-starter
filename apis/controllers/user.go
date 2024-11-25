package controllers

import (
	"gin001/apis/models"
	"gin001/services"
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
	user, err := uc.userService.CreateUser(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//
	c.JSON(http.StatusCreated, user)
}

// SignIn with user input
func (uc *UserController) SignIn(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
	})
}

// GetUser profile info.
func (uc *UserController) GetUser(c *gin.Context) {
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

// GetUsers profile info.
func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUser profile info.
func (uc *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var m models.UserInfoUpdate
	if err := c.ShouldBind(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := uc.userService.UpdateUser(uint(id), &m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
