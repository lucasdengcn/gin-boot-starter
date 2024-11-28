package controllers

import (
	"fmt"
	"gin-boot-starter/core"
	"gin-boot-starter/core/logging"
	"gin-boot-starter/infra/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ControllerBase define
type ControllerBase struct{}

func (c *ControllerBase) deferTxCallback(ctx *gin.Context, val any) {
	logging.Debug(ctx).Msgf("In recover call. Err is: %v", val)
	if val != nil {
		db.RollbackTx(ctx)
		if c.responseAsServiceError(ctx, val) {
			return
		}
		if c.responseAs404Error(ctx, val) {
			return
		}
		c.responseAs500Error(ctx, val)
	} else {
		db.CommitTx(ctx)
	}
}

func isError(val any) bool {
	if _, ok := val.(error); ok {
		return true
	} else {
		return false
	}
}

func (c *ControllerBase) responseAsServiceError(ctx *gin.Context, val any) bool {
	err, ok := val.(core.ServiceError)
	if ok {
		ctx.JSON(http.StatusInternalServerError, core.NewProblemServiceDetail(err, ctx))
		return true
	}
	return false
}

func (c *ControllerBase) responseAsRepositoryError(ctx *gin.Context, val any) bool {
	err, ok := val.(core.RepositoryError)
	if ok {
		ctx.JSON(http.StatusInternalServerError, core.NewProblemRepositoryDetail(err, ctx))
		return true
	}
	return false
}

func (c *ControllerBase) responseAs404Error(ctx *gin.Context, val any) bool {
	err, ok := val.(core.EntityNotFoundError)
	if ok {
		ctx.JSON(http.StatusNotFound, core.NewProblem404Detail(err, ctx))
		return true
	}
	return false
}

func (c *ControllerBase) responseAs500Error(ctx *gin.Context, val any) {
	err, ok := val.(error)
	if ok {
		ctx.JSON(http.StatusInternalServerError, core.NewUnexpectedDetail(err, ctx))
	} else {
		ctx.JSON(http.StatusInternalServerError, core.NewUnexpectedDetail(fmt.Errorf("Unexpected error: %v", val), ctx))
	}
}
