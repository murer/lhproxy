package util

import (
	"bytes"
	"testing"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Code uint8
	Payload [2]byte
}

func TestBinEnc(t *testing.T) {
	assert.Equal(t, []byte{0x80}, BinEnc(uint8(0x80)))
	var n uint8
	BinDec([]byte{0x80}, &n)
	assert.Equal(t, uint8(0x80), n)
}

func TestBinWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	BinWrite(buf, uint8(0x80))
	assert.Equal(t, []byte{0x80}, buf.Bytes())
	var n uint8
	BinRead(buf, &n)
	assert.Equal(t, uint8(0x80), n)
}

func TestBinStruct(t *testing.T) {
	data := &Data{7, [2]byte{1,2}}
	assert.Equal(t, []byte{7,1,2}, BinEnc(data))
	BinDec([]byte{8,2,3}, data)
	assert.Equal(t, &Data{8, [2]byte{2,3}}, data)
}
