package controller

import (
	"gin-boot-starter/api/model"
	"gin-boot-starter/core"
	"gin-boot-starter/core/exception"
	"gin-boot-starter/core/security"
	"gin-boot-starter/infra/db"
	"gin-boot-starter/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController struct
type UserController struct {
	userService *service.UserService
	aclService  *service.AclService
}

// NewUserController with dependencies services
func NewUserController(UserService *service.UserService, AclService *service.AclService) *UserController {
	return &UserController{userService: UserService, aclService: AclService}
}

// GetCurrentUser profile info.
// GetCurrentUser godoc
// @Summary Get current user information.
// @Tags UserController
// @Accept application/json
// @Produce json
// @Success 200 {object} model.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /users/v1/session [GET]
func (uc *UserController) GetCurrentUser(ctx *gin.Context) {
	//
	user := security.CurrentUser(ctx)
	ctx.JSON(http.StatusOK, user)
}

// GetUser profile info.
// GetUser godoc
// @Summary Get user summary information.
// @Tags UserController
// @Accept application/json
// @Produce json
// @Param        id   path      int  true  "Account ID"
// @Success 200 {object} model.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /users/v1/:id [GET]
func (uc *UserController) GetUser(c *gin.Context) {
	id, err := core.UintFromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.NewProblemValidationDetail("id", err.Error(), c))
		return
	}
	// check current user's ACL
	sessionUser := security.CurrentUser(c)
	uc.aclService.HasPolicy(c, sessionUser.GetID(), "user", "read")
	// query on user info
	user := uc.userService.GetUser(c, id)
	c.JSON(http.StatusOK, user)
}

// GetUsers profile info.
// @Summary Query users in paging.
// @Tags UserController
// @Accept application/json
// @Produce json
// @Param        size   path      int  true  "amount of items to return"
// @Param        page   path      int  true  "current page index"
// @Success 200 {object} model.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /users/v1/paging/:size/:page [GET]
func (uc *UserController) GetUsers(c *gin.Context) {
	users := uc.userService.GetUsers(c)
	c.JSON(http.StatusOK, users)
}

// UpdateUser profile info.
// @Summary Update user with inputs.
// @Tags UserController
// @Accept application/json
// @Produce json
// @Param model body model.UserInfoUpdate true "user info"
// @Success 200 {object} model.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /users/v1/:id [PUT]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exception.NewProblemValidationDetail("id", err.Error(), ctx))
		return
	}
	var m model.UserInfoUpdate
	if err := ctx.ShouldBind(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, exception.NewProblemBindingDetail(err, ctx))
		return
	}
	// given the request would modify data, then it should be in tx scope
	db.BeginTx(ctx)
	defer func() {
		err := recover()
		db.RecoverErrorHandle(ctx, err)
	}()
	//
	user := uc.userService.UpdateUser(ctx, uint(id), &m)
	// other service
	// other service
	// finally
	//
	ctx.JSON(http.StatusOK, user)
}
