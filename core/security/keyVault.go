package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gin001/config"
	"gin001/core/logging"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Get Public Key from public key pem file
func GetPublicKey(pemFile string) (*rsa.PublicKey, error) {
	// Read the PEM file
	data, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}

	// Parse the PEM file to extract the public key
	// Convert the public key to a string
	return ParsePublicKey(data)
}

// GetCfgPublicKey from configuration
func GetCfgPublicKey() (*rsa.PublicKey, error) {
	cfg := config.GetConfig().Security.JWT
	if strings.HasSuffix(cfg.PublicKey, ".pem") {
		path := config.GetConfig().Application.CfgPath + "/" + cfg.PublicKey
		return GetPublicKey(path)
	}
	// TO DO: Read from key vault or jwk url
	return nil, nil
}

// ParsePublicKey from pem format
func ParsePublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("Decode PEM Error")
	}
	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid PEM format: %s", block.Type)
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	return publicKey.(*rsa.PublicKey), err
}

// Get Private Key from private key pem file
func GetPrivateKey(pemFile string) (*rsa.PrivateKey, error) {
	// Read the PEM file
	data, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}

	// Parse the PEM file to extract the private key
	// Convert the private key to a string
	return ParsePrivateKey(data)
}

// ParsePrivateKey from pem format
// PKCS1 private keys typically start with -----BEGIN RSA PRIVATE KEY-----
// PKCS8 private keys typically start with -----BEGIN PRIVATE KEY-----.
func ParsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("Decode PEM Error")
	}
	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("Only support [PRIVATE KEY], invalid PEM format: %s", block.Type)
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	return privateKey.(*rsa.PrivateKey), err
}

// GetCfgPrivateKey from configuration, env
func GetCfgPrivateKey() (*rsa.PrivateKey, error) {
	jwt := config.GetConfig().Security.JWT
	if strings.HasSuffix(jwt.PrivateKey, ".pem") {
		path := config.GetConfig().Application.CfgPath + "/" + jwt.PrivateKey
		return GetPrivateKey(path)
	}
	// TO DO: Read from key vault
	return nil, nil
}

// PublicJwtKeyfuncCtx for golang-jwt/v5
func PublicJwtKeyfuncCtx(c *gin.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		cfg := config.GetConfig().Security.JWT
		if strings.HasSuffix(cfg.PublicKey, ".pem") {
			key, err := GetCfgPublicKey()
			// logging.Debug(c).Msgf("Public Key is : %T", key)
			if err != nil {
				logging.Error(c).Err(err).Msgf("Load Public Key Error")
			}
			return key, err
		}
		return nil, nil
	}
}

// SignAccessToken for sub, aud
func SignAccessToken(sub string, aud string) (string, time.Time, error) {
	jwtCfg := config.GetConfig().Security.JWT
	mySigningKey, err := GetCfgPrivateKey()
	if err != nil {
		return "", time.Time{}, err
	}
	// Create the Claims
	now := time.Now()
	expireTime := now.Add(time.Duration(jwtCfg.ExpirationTimeMinutes) * time.Minute)
	claims := &AuthClaims{
		Roles:  []string{"User"},
		Groups: []string{"User"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jwtCfg.Issuer,
			Subject:   sub,
			Audience:  []string{aud},
			ID:        jwtCfg.KeyID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	ss, err := token.SignedString(mySigningKey)
	return ss, expireTime, err
}

// SignRefreshToken for sub, aud
func SignRefreshToken(sub string, aud string) (string, time.Time, error) {
	jwtCfg := config.GetConfig().Security.JWT
	mySigningKey, err := GetCfgPrivateKey()
	if err != nil {
		return "", time.Time{}, err
	}
	// Create the Claims
	now := time.Now()
	expireTime := now.Add(time.Duration(jwtCfg.RefreshTokenExpirationTimeHours) * time.Hour)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expireTime),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    jwtCfg.Issuer,
		Subject:   sub,
		Audience:  []string{aud},
		ID:        jwtCfg.KeyID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	ss, err := token.SignedString(mySigningKey)
	return ss, expireTime, err
}
