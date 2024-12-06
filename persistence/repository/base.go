package repository

import (
	"gin-boot-starter/core/logging"
	"gin-boot-starter/infra/db"

	"github.com/gin-gonic/gin"

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

func (repo *TransactionRepo) prepareStatement(ctx *gin.Context, sql string) *sqlx.Stmt {
	tx := db.GetTx(ctx)
	if tx == nil {
		stmt, err := repo.dbCon.Preparex(sql)
		if err != nil {
			logging.Panic(ctx).Msgf("Preparex statement Error: %v, %v", sql, err)
		}
		return stmt
	}
	//
	{
		stmt, err := tx.Preparex(sql)
		if err != nil {
			logging.Panic(ctx).Msgf("Preparex statement Error: %v, %v", sql, err)
		}
		return stmt
	}
}

func (repo *TransactionRepo) prepareNamed(ctx *gin.Context, sql string) *sqlx.NamedStmt {
	tx := db.GetTx(ctx)
	if tx == nil {
		stmt, err := repo.dbCon.PrepareNamed(sql)
		if err != nil {
			logging.Panic(ctx).Msgf("PrepareNamed statement Error: %v, %v", sql, err)
		}
		return stmt
	}
	//
	{
		logging.Debug(ctx).Msgf("Exec in tx %v", tx)
		stmt, err := tx.PrepareNamed(sql)
		if err != nil {
			logging.Panic(ctx).Msgf("PrepareNamed statement Error: %v, %v", sql, err)
		}
		return stmt
	}
}
