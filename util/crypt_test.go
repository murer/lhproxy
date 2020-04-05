package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPKCS5(t *testing.T) {
	plaintext := []byte{5, 6, 7}
	padded := []byte{5, 6, 7, 5, 5, 5, 5, 5}
	assert.Equal(t, padded, pkcs5pad(plaintext, 8))
	assert.Equal(t, plaintext, pkcs5trim(padded, 8))
}

func TestEncrypt(t *testing.T) {
	plaintext := []byte{5, 6, 7}
	cryptor := &Cryptor{}
	assert.Equal(t, CRYPTOR_KEY_SIZE, len(cryptor.GenSecret()))
	c1 := cryptor.Encrypt(plaintext)
	c2 := cryptor.Encrypt(plaintext)
	assert.NotEqual(t, plaintext, c1)
	assert.NotEqual(t, plaintext, c2)
	assert.NotEqual(t, c2, c1)
	assert.Equal(t, plaintext, cryptor.Decrypt(c1))
	assert.Equal(t, plaintext, cryptor.Decrypt(c2))

	assert.Equal(t, 32, len(c1))
	assert.Equal(t, 32, len(c2))
}
