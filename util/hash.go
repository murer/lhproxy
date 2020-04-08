package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(data []byte) []byte {
	ret := sha256.Sum256(data)
	return ret[:]
}

func SHA256Hex(data []byte) string {
	ret := SHA256(data)
	return hex.EncodeToString(ret)
}
