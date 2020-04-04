package util

import (
	"crypto/sha256"
)

func SHA256(data []byte) []byte {
	ret := sha256.Sum256(data)
	return ret[:]
}
