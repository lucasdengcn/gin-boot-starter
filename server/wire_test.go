package server

import (
	"gin001/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBWire(t *testing.T) {
	var basePath = config.GetBasePath()
	config.LoadConf(basePath, "test")
	uc := InitializeUserController()
	assert.NotNil(t, uc)
}
