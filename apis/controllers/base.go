package controllers

import (
	"gin-boot-starter/core"
	"gin-boot-starter/core/logging"
	"gin-boot-starter/infra/db"

	"github.com/gin-gonic/gin"
)

// ControllerBase define
type ControllerBase struct{}

func (c *ControllerBase) deferTxCallback(ctx *gin.Context, val any) {
	logging.Debug(ctx).Msgf("In recover call. Err is: %v", val)
	if val != nil {
		db.RollbackTx(ctx)
		core.ResponseOnError(ctx, val)
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
