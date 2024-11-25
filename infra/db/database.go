package db

import (
	"context"
	"fmt"
	"gin001/config"
	"log"

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
	dbCon, err := sqlx.Open("pgx", cfg.GetString("db.url"))
	if err != nil {
		return nil, err
	}
	dbCon.SetMaxIdleConns(cfg.GetInt("db.pool.min"))
	dbCon.SetMaxOpenConns(cfg.GetInt("db.pool.max"))
	fmt.Println("DB Connect Successfully.")
	err = dbCon.Ping()
	if err != nil {
		return nil, err
	}
	dbSQL = dbCon
	return dbCon, nil
}

// GetDBCon provider return db instance
func GetDBCon() *sqlx.DB {
	if dbSQL == nil {
		dbCon, err := ConnectDB()
		if err != nil {
			log.Fatal("DB Connect Failed.")
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
}

// BeginTx return context
func BeginTx(ctx context.Context) context.Context {
	if dbSQL != nil {
		tx, err := dbSQL.BeginTxx(ctx, nil)
		if err != nil {
			log.Panicln("Begin Tx Error", err)
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
		log.Printf("tx Commit Error. %v\n", dbTx)
		panic(err)
	}
	log.Printf("tx Commit Success. %v\n", dbTx)
	dbTx = nil
	return context.WithValue(ctx, txKey, nil)
}

// RollbackTx return error
func RollbackTx(ctx context.Context) {
	dbTx := GetTx(ctx)
	err := dbTx.Rollback()
	if err != nil {
		log.Panicf("tx Rollback Error. %v\n", dbTx)
	} else {
		log.Printf("tx Rollback Success. %v\n", dbTx)
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
		log.Panicf("Can't Convert Tx object from context. %v", dbTx)
	}
	return dbTx
}
