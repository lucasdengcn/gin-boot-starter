package controllers

import (
	"context"
	"gin001/apis/models"
	"gin001/infra/db"
	"gin001/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController struct
type UserController struct {
	ControllerBase
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
	ctx := db.BeginTx(context.Background())
	defer func() {
		err := recover()
		uc.deferTxCallback(ctx, c, err)
	}()
	//
	user, err := uc.userService.CreateUser(ctx, &m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// other service
	// other service
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
	ctx := context.Background()
	user, err := uc.userService.GetUser(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUsers profile info.
func (uc *UserController) GetUsers(c *gin.Context) {
	ctx := context.Background()
	users, err := uc.userService.GetUsers(ctx)
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
	// given the request would modify data, then it should be in tx scope
	ctx := db.BeginTx(context.Background())
	defer func() {
		err := recover()
		uc.deferTxCallback(ctx, c, err)
	}()
	//
	user, err := uc.userService.UpdateUser(ctx, uint(id), &m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// other service
	// other service
	// finally
	//
	c.JSON(http.StatusOK, user)
}
