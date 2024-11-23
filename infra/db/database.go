package db

import (
	"fmt"
	"gin001/config"
	"log"

	// postgresql driver
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var dbSQL *sqlx.DB

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
