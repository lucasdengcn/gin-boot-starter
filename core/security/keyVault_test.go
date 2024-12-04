package security

import (
	"gin-boot-starter/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.LoadConf(config.GetBasePath(), "test")
}

func TestGetCfgPublicKey(t *testing.T) {
	pk, err := GetCfgPublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pk)
}

func TestGetCfgPrivateKey(t *testing.T) {
	pk, err := GetCfgPrivateKey()
	assert.NoError(t, err)
	assert.NotNil(t, pk)
}

func TestSignAccessToken(t *testing.T) {
	token, exp, err := SignAccessToken("1", "test")
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotNil(t, exp)
	assert.True(t, exp.After(time.Now()))
}

func TestSignRefreshToken(t *testing.T) {
	token, exp, err := SignRefreshToken("1", "test")
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotNil(t, exp)
	assert.True(t, exp.After(time.Now()))
}

func TestParseAccessToken(t *testing.T) {
	tokenString, exp, err := SignAccessToken("1", "test")
	assert.NoError(t, err)
	assert.NotNil(t, tokenString)
	assert.NotNil(t, exp)
	assert.True(t, exp.After(time.Now()))
	//
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, PublicJwtKeyfuncCtx(nil))
	assert.NoError(t, err)
	assert.NotNil(t, token)
}
