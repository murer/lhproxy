package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	plaintext := []byte{5, 6, 7}
	c1 := Encrypt(plaintext)
	c2 := Encrypt(plaintext)
	assert.NotEqual(t, plaintext, c1)
	assert.NotEqual(t, plaintext, c2)
	assert.NotEqual(t, c2, c1)
	assert.Equal(t, plaintext, Decrypt(c1))
	assert.Equal(t, plaintext, Decrypt(c2))
}
