package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConf(t *testing.T) {
	var basePath = GetBasePath()
	err := LoadConf(basePath, "test")
	assert.NoError(t, err)
}

func TestLoadConfFail(t *testing.T) {
	var basePath = GetBasePath()
	err := LoadConf(basePath, "non-existent-file")
	assert.Error(t, err)
}

func TestConfValue(t *testing.T) {
	var basePath = GetBasePath()
	LoadConf(basePath, "test")
	cfg := GetConfig()
	assert.Equal(t, "Example test", cfg.Application.Name)
}
