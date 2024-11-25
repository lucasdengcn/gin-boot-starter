package controllers

import (
	"context"
	"gin001/infra/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ControllerBase define
type ControllerBase struct{}

func (c *ControllerBase) deferTxCallback(ctx context.Context, ginContext *gin.Context, err any) {
	log.Printf("In recover call. Err is: %v\n", err)
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
