package repository

import (
	"context"
	"gin001/infra/db"

	"github.com/rs/zerolog/log"

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
			log.Panic().Msgf("Preparex statement Error: %v, %v", sql, err)
		}
		return stmt
	}
	//
	{
		stmt, err := tx.Preparex(sql)
		if err != nil {
			log.Panic().Msgf("Preparex statement Error: %v, %v", sql, err)
		}
		return stmt
	}
}

func (repo *TransactionRepo) prepareNamed(ctx context.Context, sql string) *sqlx.NamedStmt {
	tx := db.GetTx(ctx)
	if tx == nil {
		stmt, err := repo.dbCon.PrepareNamed(sql)
		if err != nil {
			log.Panic().Msgf("PrepareNamed statement Error: %v, %v", sql, err)
		}
		return stmt
	}
	//
	{
		log.Debug().Msgf("Exec in tx %v", tx)
		stmt, err := tx.PrepareNamed(sql)
		if err != nil {
			log.Panic().Msgf("PrepareNamed statement Error: %v, %v", sql, err)
		}
		return stmt
	}
}
