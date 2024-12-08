package db

import (
	"gin-boot-starter/core/config"
	"gin-boot-starter/core/exception"
	"gin-boot-starter/core/logging"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	// postgresql driver
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var once sync.Once

var dbSQL *sqlx.DB

const txKey string = "txScoped"

// ConnectDB db connections
func ConnectDB() (*sqlx.DB, error) {
	if dbSQL != nil {
		return dbSQL, nil
	}
	var err error = nil
	once.Do(func() {
		cfg := config.GetConfig()
		dbCon, err := sqlx.Open(cfg.DataSource.Driver, cfg.DataSource.URL)
		if err != nil {
			return
		}
		dbCon.SetMaxIdleConns(cfg.DataSource.PoolMin)
		dbCon.SetMaxOpenConns(cfg.DataSource.PoolMax)
		err = dbCon.Ping()
		if err != nil {
			log.Fatal().Err(err).Msg("DB Ping Failed.")
			return
		}
		dbSQL = dbCon
		log.Info().Msg("DB Connect Successfully.")
	})
	return dbSQL, err
}

// GetDBCon provider return db instance
func GetDBCon() *sqlx.DB {
	if dbSQL == nil {
		dbCon, err := ConnectDB()
		if err != nil {
			log.Fatal().Msg("DB Connect Failed.")
		}
		return dbCon
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
func BeginTx(ctx *gin.Context) {
	if dbSQL == nil {
		panic("DB not init yet.")
	}
	tx, err := dbSQL.BeginTxx(ctx, nil)
	if err != nil {
		logging.Panic(ctx).Msgf("Begin Tx Error: %v", err)
	}
	ctx.Set(txKey, tx)
}

// CommitTx return context
func CommitTx(ctx *gin.Context) {
	dbTx := GetTx(ctx)
	if dbTx == nil {
		logging.Error(ctx).Msgf("tx Commit, but No Tx attached to the context. check the caller chain")
		return
	}
	err := dbTx.Commit()
	if err != nil {
		logging.Error(ctx).Msgf("tx Commit Error. %v", dbTx)
		panic(err)
	}
	logging.Debug(ctx).Msgf("tx Commit Success. %v", dbTx)
	dbTx = nil
	ctx.Set(txKey, nil)
}

// RollbackTx return error
func RollbackTx(ctx *gin.Context) {
	dbTx := GetTx(ctx)
	if dbTx == nil {
		logging.Error(ctx).Msgf("tx Rollback, but No Tx attached to the context. check the caller chain")
		return
	}
	err := dbTx.Rollback()
	if err != nil {
		logging.Error(ctx).Msgf("tx Rollback Error. %v", dbTx)
	} else {
		logging.Debug(ctx).Msgf("tx Rollback Success. %v", dbTx)
	}
}

// GetTx return sqlx.Tx
func GetTx(ctx *gin.Context) *sqlx.Tx {
	val, exists := ctx.Get(txKey)
	if val == nil || !exists {
		return nil
	}
	dbTx, ok := val.(*sqlx.Tx)
	if !ok {
		logging.Panic(ctx).Msgf("Can't Convert Tx object from context. %v", dbTx)
	}
	return dbTx
}

// RecoverErrorHandle, to recover from panic.
func RecoverErrorHandle(ctx *gin.Context, r any) {
	if r != nil {
		RollbackTx(ctx)
		exception.ResponseOnError(ctx, r)
	} else {
		CommitTx(ctx)
	}
}
