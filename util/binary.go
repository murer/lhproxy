package util

import (
	"bytes"
	"encoding/binary"
)

func NewBinary(data []byte) *Binary {
	return &Binary{
		buf: bytes.NewBuffer(data),
	}
}

type Binary struct {
	buf *bytes.Buffer
}

func (b *Binary) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *Binary) WriteBytes(s []byte) {
	b.WriteUInt16(uint16(len(s)))
	_, err := b.buf.Write(s)
	Check(err)
}

func (b *Binary) ReadBytes() []byte {
	l := b.ReadUInt16()
	return ReadFully(b.buf, int(l))
}

func (b *Binary) WriteString(s string) {
	b.WriteBytes([]byte(s))
}

func (b *Binary) ReadString() string {
	return string(b.ReadBytes())
}

func (b *Binary) WriteUInt16(n uint16) {
	Check(binary.Write(b.buf, binary.BigEndian, n))
}

func (b *Binary) ReadUInt16() uint16 {
	var ret uint16
	Check(binary.Read(b.buf, binary.BigEndian, &ret))
	return ret
}
