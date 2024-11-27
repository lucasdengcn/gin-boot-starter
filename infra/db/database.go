package db

import (
	"context"
	"gin001/config"

	"github.com/rs/zerolog/log"

	// postgresql driver
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var dbSQL *sqlx.DB

type txKeyType string

const txKey txKeyType = "txScoped"

// ConnectDB db connections
func ConnectDB() (*sqlx.DB, error) {
	if dbSQL != nil {
		return dbSQL, nil
	}
	cfg := config.GetConfig()
	dbCon, err := sqlx.Open(cfg.DataSource.Driver, cfg.DataSource.URL)
	if err != nil {
		return nil, err
	}
	dbCon.SetMaxIdleConns(cfg.DataSource.PoolMin)
	dbCon.SetMaxOpenConns(cfg.DataSource.PoolMax)
	err = dbCon.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("DB Ping Failed.")
		return nil, err
	}
	dbSQL = dbCon
	log.Info().Msg("DB Connect Successfully.")
	return dbCon, nil
}

// GetDBCon provider return db instance
func GetDBCon() *sqlx.DB {
	if dbSQL == nil {
		dbCon, err := ConnectDB()
		if err != nil {
			log.Fatal().Msg("DB Connect Failed.")
		}
		dbCon.Ping()
		dbSQL = dbCon
	}
	return dbSQL
}

// Close DB connection
func Close() {
	if dbSQL != nil {
		dbSQL.Close()
	}
	log.Info().Msg("DB Shutdown.")
}

// BeginTx return context
func BeginTx(ctx context.Context) context.Context {
	if dbSQL != nil {
		tx, err := dbSQL.BeginTxx(ctx, nil)
		if err != nil {
			log.Panic().Msgf("Begin Tx Error: %v", err)
		}
		ctxNew := context.WithValue(ctx, txKey, tx)
		return ctxNew
	}
	panic("DB not init yet.")
}

// CommitTx return context
func CommitTx(ctx context.Context) context.Context {
	dbTx := GetTx(ctx)
	err := dbTx.Commit()
	if err != nil {
		log.Error().Msgf("tx Commit Error. %v", dbTx)
		panic(err)
	}
	log.Debug().Msgf("tx Commit Success. %v", dbTx)
	dbTx = nil
	return context.WithValue(ctx, txKey, nil)
}

// RollbackTx return error
func RollbackTx(ctx context.Context) {
	dbTx := GetTx(ctx)
	err := dbTx.Rollback()
	if err != nil {
		log.Error().Msgf("tx Rollback Error. %v", dbTx)
	} else {
		log.Debug().Msgf("tx Rollback Success. %v", dbTx)
	}
}

// GetTx return sqlx.Tx
func GetTx(ctx context.Context) *sqlx.Tx {
	val := ctx.Value(txKey)
	if val == nil {
		return nil
	}
	dbTx, ok := val.(*sqlx.Tx)
	if !ok {
		log.Panic().Msgf("Can't Convert Tx object from context. %v", dbTx)
	}
	return dbTx
}
