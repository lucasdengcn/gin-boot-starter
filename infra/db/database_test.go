package db

import (
	"fmt"
	"gin-boot-starter/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConnectionString(t *testing.T) {
	var basePath = config.GetBasePath()
	config.LoadConf(basePath, "test")
	cfg := config.GetConfig()
	assert.NotEmpty(t, cfg.DataSource.URL)
}

func TestDBConnection(t *testing.T) {
	// Arrange
	var basePath = config.GetBasePath()
	config.LoadConf(basePath, "test")
	db, err := ConnectDB()
	assert.NoError(t, err)
	//
	rows, err := db.Query("select NOW()")
	assert.NoError(t, err)
	fmt.Println(rows)
	db.Close()
}
