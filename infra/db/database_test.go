package db

import (
	"fmt"
	"gin001/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConnectionString(t *testing.T) {
	config.LoadConf("test")
	cfg := config.GetConfig()
	assert.NotEmpty(t, cfg.GetString("db.url"))
}

func TestDBConnection(t *testing.T) {
	// Arrange
	config.LoadConf("test")
	db, err := ConnectDB()
	assert.NoError(t, err)
	//
	rows, err := db.Query("select NOW()")
	assert.NoError(t, err)
	fmt.Println(rows)
	db.Close()
}
