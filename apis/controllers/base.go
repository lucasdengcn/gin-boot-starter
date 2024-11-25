package controllers

import (
	"context"
	"gin001/infra/db"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

// ControllerBase define
type ControllerBase struct{}

func (c *ControllerBase) deferTxCallback(ctx context.Context, ginContext *gin.Context, err any) {
	log.Debug().Msgf("In recover call. Err is: %v", err)
	if err != nil {
		db.RollbackTx(ctx)
		realErr, ok := err.(error)
		if ok {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": realErr.Error()})
		} else {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected Error"})
		}
	} else {
		db.CommitTx(ctx)
	}
}
