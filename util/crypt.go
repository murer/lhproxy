package util

import (
	"bytes"
	"log"
	"crypto/rand"
	"crypto/aes"
	"crypto/cipher"
)

const CRYPTOR_BLOCK_SIZE = 16

func pkcs5pad(plaintext []byte, blockSize int) []byte {
	plaintextLen := len(plaintext)
	padLen := blockSize - (plaintextLen % blockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(plaintext, padText...)
}

func pkcs5trim(ciphertext []byte, blockSize int) []byte {
	ciphertextLen := len(ciphertext)
	paddingLen := int(ciphertext[ciphertextLen-1])
	if paddingLen >= ciphertextLen || paddingLen > blockSize {
		log.Panicf("Wrong padding. blockSize: %d, paddingLen: %d, ciphertextLen: %d", blockSize, paddingLen, ciphertextLen)
	}
	return ciphertext[:ciphertextLen-paddingLen]
}

type Cryptor struct {
	Secret []byte
}

func (c *Cryptor) BlockSize() int {
	if c.Secret == nil || len(c.Secret) == 0 {
		log.Panic("Cryptor is not ready")
	}
	return len(c.Secret)
}

func (c *Cryptor) Gen() []byte {
	key := make([]byte, CRYPTOR_BLOCK_SIZE)
	n, err := rand.Read(key)
	Check(err)
	if n != CRYPTOR_BLOCK_SIZE {
		log.Panicf("wrong: %d, expected: %d", n, CRYPTOR_BLOCK_SIZE)
	}
	return key
}

func (c *Cryptor) GenSecret() []byte {
	c.Secret = c.Gen()
	return c.Secret
}

func (c *Cryptor) Encrypt(plaintext []byte) []byte {
	block, err := aes.NewCipher(c.Secret)
	Check(err)
	iv := []byte("1234567890123456")
	log.Printf("xxx %d - %d", c.BlockSize(), len(iv))
	encrypter := cipher.NewCBCEncrypter(block, iv)
	padded := pkcs5pad(plaintext, c.BlockSize())
	encrypter.CryptBlocks(padded, padded)
	return padded
}

func (c *Cryptor) Decrypt(ciphertext []byte) []byte {
	block, err := aes.NewCipher(c.Secret)
	Check(err)
	iv := []byte("1234567890123456")
	log.Printf("xxx %d - %d", c.BlockSize(), len(iv))
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(ciphertext, ciphertext)
	trimmed := pkcs5trim(ciphertext, c.BlockSize())
	return trimmed
}
