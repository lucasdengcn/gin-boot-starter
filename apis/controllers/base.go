package controllers

import (
	"context"
	"gin001/infra/db"
	"log"
)

// ControllerBase define
type ControllerBase struct{}

func (c *ControllerBase) deferCallback(ctx context.Context) {
	// log.Println("In defer call")
	if err := recover(); err != nil {
		db.RollbackTx(ctx)
		log.Println(err)
	} else {
		db.CommitTx(ctx)
	}
}
