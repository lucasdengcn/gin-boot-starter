package server

import (
	"gin-boot-starter/core/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func before() {
	var basePath = config.GetBasePath()
	config.LoadConf(basePath, "test")
}

func after() {

}

func TestMain(m *testing.M) {
	before()
	code := m.Run()
	after()
	os.Exit(code)
}

func TestInitializeUserController(t *testing.T) {
	uc := InitializeUserController()
	assert.NotNil(t, uc)
}

func TestInitializeAccountController(t *testing.T) {
	uc := InitializeAccountController()
	assert.NotNil(t, uc)
}
