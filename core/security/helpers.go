package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

var emptyString = ""

// AesEncrypt to encrypt plain text with AES algorithm.
// key length MUST BE 32 characters
func AesEncrypt(plainText string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return emptyString, err
	}
	bs := []byte(plainText)
	cipherText := make([]byte, aes.BlockSize+len(bs))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return emptyString, err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], bs)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AesDecrypt to decrypt encrypted text with AES algorithm.
func AesDecrypt(cipherText string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return emptyString, err
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return emptyString, err
	}
	iv := cipherBytes[:aes.BlockSize]
	cipherTextBytes := cipherBytes[aes.BlockSize:]
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherTextBytes, cipherTextBytes)
	return string(cipherTextBytes), nil
}

// HashPassword to hash plain text password for store
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return emptyString, err
	}
	return string(hash), nil
}

// VerifyPassword to verify input with hashed-password
func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RSAEncrypt to encrypt plain text using RSA with public key
func RSAEncrypt(plainText []byte, publicKey []byte) (string, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return emptyString, fmt.Errorf("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return emptyString, err
	}
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pub, plainText)
	if err != nil {
		return emptyString, err
	}
	return string(cipherText), nil
}

// RSADecrypt to decrypt encrypted text using RSA with private key
func RSADecrypt(cipherText string, privateKey []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return emptyString, fmt.Errorf("failed to decode PEM block containing private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return emptyString, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, priv, []byte(cipherText))
	if err != nil {
		return emptyString, err
	}
	return string(plainText), nil
}
