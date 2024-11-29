package controllers

import (
	"gin-boot-starter/apis/models"
	"gin-boot-starter/core"
	"gin-boot-starter/core/logging"
	"gin-boot-starter/core/security"
	"gin-boot-starter/infra/db"
	"gin-boot-starter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	ControllerBase
	userService *services.UserService
	aclService  *services.AclService
}

func NewAccountController(userService *services.UserService, aclService *services.AclService) *AccountController {
	return &AccountController{
		userService: userService,
		aclService:  aclService,
	}
}

// SignUp with user input
// SignUp godoc
// @Summary Create account on user demand.
// @Tags AccountController
// @Accept application/json
// @Produce json
// @Param model body models.UserSignUp true "user input"
// @Success 201 {object} models.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /accounts/v1/signup [POST]
func (c *AccountController) SignUp(ctx *gin.Context) {
	var m models.UserSignUp
	if err := ctx.ShouldBind(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, core.NewProblemBindingDetail(err, ctx))
		return
	}
	db.BeginTx(ctx)
	defer func() {
		err := recover()
		c.deferTxCallback(ctx, err)
	}()
	//
	user := c.userService.CreateUser(ctx, &m)
	c.aclService.SetForNewUser(ctx, user.ID)
	//
	ctx.JSON(http.StatusCreated, user)
}

// SignIn with user input
// SignIn godoc
// @Summary SignIn user on user demand.
// @Tags AccountController
// @Accept application/json
// @Produce json
// @Param model body models.UserSignIn true "user input"
// @Success 200 {object} models.AuthTokens
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /accounts/v1/signin [POST]
func (c *AccountController) SignIn(ctx *gin.Context) {
	//
	var m models.UserSignIn
	if err := ctx.ShouldBind(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, core.NewProblemBindingDetail(err, ctx))
		return
	}
	userInfo := c.userService.VerifyPassword(ctx, &m)
	//
	uid := core.StringFromUint(userInfo.ID)
	accessToken, expireTime, err := security.SignAccessToken(uid, "web")
	if err != nil {
		logging.Error(ctx).Err(err).Msg("signing access token error")
	}
	refreshToken, expireTimeRT, err := security.SignRefreshToken(uid, "web")
	if err != nil {
		logging.Error(ctx).Err(err).Msg("signing refresh token error")
	}
	//
	ctx.JSON(http.StatusCreated, models.AuthTokens{
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		ExpireIn:          &expireTime,
		RefreshTokenExpIn: &expireTimeRT,
	})
}
