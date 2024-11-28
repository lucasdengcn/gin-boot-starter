package security

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesEncryptDecrypt(t *testing.T) {
	text := "Hello, World!"
	key := "0123456789abcdef"
	cipherText, err := AesEncrypt(text, key)
	if err != nil {
		t.Error("Encryption error:", err)
		return
	}
	t.Logf("Cipher text: %v", cipherText)
	plainText, err := AesDecrypt(cipherText, key)
	if err != nil {
		t.Error("Decryption error:", err)
		return
	}
	t.Logf("Plain text: %v", plainText)
	//
	assert.Equal(t, text, plainText)
}

func TestPasswordVerification(t *testing.T) {
	text := "HelloWorld"
	cipherText, err := HashPassword(text)
	if err != nil {
		t.Error("Hash error:", err)
		return
	}
	t.Logf("Hashed text: %v", cipherText)
	ok := VerifyPassword(text, cipherText)
	assert.True(t, ok)
}

func TestRSAEncryptDecrypt(t *testing.T) {
	// RSASSA-PKCS1-v1_5 SHA-512 PEM
	text := "Hello, World!"
	privateKey := []byte(`
-----BEGIN PRIVATE KEY-----
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
F3aCdYEvFS6RCuL9lQxDSUwkmMxxiKU=
-----END PRIVATE KEY-----
  `)
	publicKey := []byte(`
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2QJ2fIrcz7pB+bQneBw/
UMJKveItKuV3vSWyzNtfoWMHkUeaZ7poq72SYv2Fxz/GiKQH0HO/9e3mp1nYgJnI
nq6qoPv2ev9VNgIy6SleMTECGw60Wpy2U1V13dUPzDNucEdwpjC7v82BXo9he/5v
j+2ZjlFs/CBNvQdnI9VzxCoObDqfnggaJ+Qtd68ULeYE6zV+OmTv9JzSPkH7F8R0
O9uxRg3JMLZ4INoI9vI7gq4Ubw711PSP8hUm69Kw5OgJsMskVr1a9LhFVThOyLC1
Vs5cMtEeFBTscJnq6yLNUGmTS9IeQ4ls/gXdqbNFSQFh0A+lnZDiH1jC+XCtQSBH
HkvY+w9sjIsFDDBwQCvwGxXc5HKcvYJCBpP0IwTMNzV9ny+u6Y+StE21Gl+bCGnT
/g1Dgu6tXNDRj6sh9PJmldl+fKOPzGxWg6f8uf2++nVfuEvfK1eChpHJLtGHvwv/
YPSQftL0O1a5h2iomsxSSu/Ev0qIcXyHsH/GAFepD1P+Hi5wVaganGzn1ReGrA30
NwTuuBdbp14gfOX63vUYQW3fFD+zLTdsIeHtRbmnUlq6QJmj69QrxbGGB9bo1QYl
/LzXYcutuKVgN44L0XWZj1NjC6KTfnFfdwkitHPTV7pLeUVQQGs5HoFKu70PQ11O
kbDKobNAKdq9YzFY+j1/DfsCAwEAAQ==
-----END PUBLIC KEY-----
  `)
	cipherText, err := RSAEncrypt([]byte(text), publicKey)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}
	fmt.Println("Cipher text:", cipherText)
	plainText, err := RSADecrypt(cipherText, privateKey)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
	fmt.Println("Plain text:", plainText)
	//
	assert.Equal(t, text, plainText)
}
