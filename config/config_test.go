package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConf(t *testing.T) {
	err := LoadConf("test")
	assert.NoError(t, err)
}

func TestLoadConfFail(t *testing.T) {
	err := LoadConf("non-existent-file")
	assert.Error(t, err)
}

func TestConfValue(t *testing.T) {
	LoadConf("test")
	cfg := GetConfig()
	assert.Equal(t, "Example test", cfg.GetString("app.name"))
}