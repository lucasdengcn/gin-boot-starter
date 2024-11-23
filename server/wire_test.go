package server

import (
	"gin001/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBWire(t *testing.T) {
	config.LoadConf("test")
	uc := InitializeUserController()
	assert.NotNil(t, uc)
	uc.GetUser(nil)
}
