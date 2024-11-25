package repository

import (
	"context"
	"gin001/infra/db"
	"log"

	"github.com/jmoiron/sqlx"
)

// TransactionRepo define
type TransactionRepo struct {
	dbCon *sqlx.DB
}

// NewTransactionRepo with dbCon
func NewTransactionRepo(dbCon *sqlx.DB) TransactionRepo {
	return TransactionRepo{
		dbCon: dbCon,
	}
}

func (repo *TransactionRepo) prepareStatement(ctx context.Context, sql string) *sqlx.Stmt {
	tx := db.GetTx(ctx)
	if tx == nil {
		stmt, err := repo.dbCon.Preparex(sql)
		if err != nil {
			log.Panicf("Preparex statement Error: %v, %v\n", sql, err)
		}
		return stmt
	}
	//
	{
		stmt, err := tx.Preparex(sql)
		if err != nil {
			log.Panicf("Preparex statement Error: %v, %v\n", sql, err)
		}
		return stmt
	}
}

func (repo *TransactionRepo) prepareNamed(ctx context.Context, sql string) *sqlx.NamedStmt {
	tx := db.GetTx(ctx)
	// log.Printf("tx %v\n", tx)
	if tx == nil {
		stmt, err := repo.dbCon.PrepareNamed(sql)
		if err != nil {
			log.Panicf("PrepareNamed statement Error: %v, %v\n", sql, err)
		}
		return stmt
	}
	//
	{
		stmt, err := tx.PrepareNamed(sql)
		if err != nil {
			log.Panicf("PrepareNamed statement Error: %v, %v\n", sql, err)
		}
		return stmt
	}
}
