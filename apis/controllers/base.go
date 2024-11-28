package controllers

import (
	"gin001/infra/db"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

// ControllerBase define
type ControllerBase struct{}

func (c *ControllerBase) deferTxCallback(ctx *gin.Context, err any) {
	log.Debug().Msgf("In recover call. Err is: %v", err)
	if err != nil {
		db.RollbackTx(ctx)
		realErr, ok := err.(error)
		if ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": realErr.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected Error"})
		}
	} else {
		db.CommitTx(ctx)
	}
}
