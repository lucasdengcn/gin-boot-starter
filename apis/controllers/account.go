package controllers

import (
	"gin001/apis/models"
	"gin001/core"
	"gin001/core/logging"
	"gin001/core/security"
	"gin001/infra/db"
	"gin001/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	ControllerBase
	userService *services.UserService
}

func NewAccountController(userService *services.UserService) *AccountController {
	return &AccountController{
		userService: userService,
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
	// other service
	// other service
	//
	ctx.JSON(http.StatusCreated, user)
}

// SignIn with user input
// SignIn godoc
// @Summary SignIn user on user demand.
// @Tags AccountController
// @Accept application/json
// @Produce json
// @Param model body models.UserSignUp true "user input"
// @Success 200 {object} models.UserInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /accounts/v1/signin [POST]
func (c *AccountController) SignIn(ctx *gin.Context) {
	//
	accessToken, expireTime, err := security.SignAccessToken("1", "web")
	if err != nil {
		logging.Error(ctx).Err(err).Msg("signing access token error")
	}
	refreshToken, expireTimeRT, err := security.SignRefreshToken("1", "web")
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
