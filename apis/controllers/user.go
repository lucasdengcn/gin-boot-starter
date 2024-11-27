package controllers

import (
	"context"
	"gin001/apis/models"
	"gin001/core"
	"gin001/infra/db"
	"gin001/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
// SignUp godoc
// @Summary Create account on user demand.
// @Tags UserController
// @Accept application/json
// @Produce json
// @Param model body models.UserSignUp true "user input"
// @Success 201 {object} models.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /users/v1/signup [POST]
func (uc *UserController) SignUp(c *gin.Context) {
	var m models.UserSignUp
	if err := c.ShouldBind(&m); err != nil {
		c.JSON(http.StatusBadRequest, core.NewBindingError(err, c))
		return
	}
	ctx := db.BeginTx(context.Background())
	defer func() {
		err := recover()
		uc.deferTxCallback(ctx, c, err)
	}()
	//
	user := uc.userService.CreateUser(ctx, &m)
	// other service
	// other service
	//
	c.JSON(http.StatusCreated, user)
}

// SignIn with user input
// SignIn godoc
// @Summary SignIn user on user demand.
// @Tags UserController
// @Accept application/json
// @Produce json
// @Param model body models.UserSignUp true "user input"
// @Success 200 {object} models.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /users/v1/signin [POST]
func (uc *UserController) SignIn(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
	})
}

// GetUser profile info.
// GetUser godoc
// @Summary Get user summary information.
// @Tags UserController
// @Accept application/json
// @Produce json
// @Param        id   path      int  true  "Account ID"
// @Success 200 {object} models.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /users/v1/:id [GET]
func (uc *UserController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, core.NewValidationError("id", err.Error(), c))
		return
	}
	log.Debug().Msgf("GetUser with id:%v", id)
	ctx := context.Background()
	user := uc.userService.GetUser(ctx, uint(id))
	c.JSON(http.StatusOK, user)
}

// GetUsers profile info.
// @Summary Query users in paging.
// @Tags UserController
// @Accept application/json
// @Produce json
// @Param        size   path      int  true  "amount of items to return"
// @Param        page   path      int  true  "current page index"
// @Success 200 {object} models.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /users/v1/paging/:size/:page [GET]
func (uc *UserController) GetUsers(c *gin.Context) {
	ctx := context.Background()
	users := uc.userService.GetUsers(ctx)
	c.JSON(http.StatusOK, users)
}

// UpdateUser profile info.
// @Summary Update user with inputs.
// @Tags UserController
// @Accept application/json
// @Produce json
// @Param model body models.UserInfoUpdate true "user info"
// @Success 200 {object} models.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /users/v1/:id [PUT]
func (uc *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, core.NewValidationError("id", err.Error(), c))
		return
	}
	var m models.UserInfoUpdate
	if err := c.ShouldBind(&m); err != nil {
		c.JSON(http.StatusBadRequest, core.NewBindingError(err, c))
		return
	}
	// given the request would modify data, then it should be in tx scope
	ctx := db.BeginTx(context.Background())
	defer func() {
		err := recover()
		uc.deferTxCallback(ctx, c, err)
	}()
	//
	user := uc.userService.UpdateUser(ctx, uint(id), &m)
	// other service
	// other service
	// finally
	//
	c.JSON(http.StatusOK, user)
}
