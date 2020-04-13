package util

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
)

func NewBinary(data []byte) *Binary {
	return NewBinaryBuffer(bytes.NewBuffer(data))
}

func NewBinaryBuffer(buf *bytes.Buffer) *Binary {
	return &Binary{buffer: buf, reader: buf, writer: buf}
}

func NewBinaryReader(r io.Reader) *Binary {
	return &Binary{reader: r}
}

func NewBinaryWriter(w io.Writer) *Binary {
	return &Binary{writer: w}
}

func NewBinaryReadWriter(rw io.ReadWriter) *Binary {
	return &Binary{reader: rw, writer: rw}
}

type Binary struct {
	buffer *bytes.Buffer
	reader io.Reader
	writer io.Writer
}

func (b *Binary) Bytes() []byte {
	return b.buffer.Bytes()
}

func (b *Binary) WriteNillableBytes(s []byte) {
	if s == nil {
		b.WriteUInt8(0)
		return
	}
	b.WriteUInt8(1)
	b.WriteBytes(s)
}

func (b *Binary) WriteBytes(s []byte) {
	if s == nil {
		log.Panic("Use WriteNillableBytes to support nil")
	}
	b.WriteUInt16(uint16(len(s)))
	_, err := b.writer.Write(s)
	Check(err)
}

func (b *Binary) ReadNillableBytes() []byte {
	n := b.ReadUInt8()
	if n == 0 {
		return nil
	}
	if n != 1 {
		log.Panicf("non nil []bytes must be 0x01, but was: %x", n)
	}
	return b.ReadBytes()
}

func (b *Binary) ReadBytes() []byte {
	l := b.ReadUInt16()
	return ReadFully(b.reader, int(l))
}

func (b *Binary) WriteString(s string) {
	b.WriteBytes([]byte(s))
}

func (b *Binary) ReadString() string {
	return string(b.ReadBytes())
}

func (b *Binary) WriteUInt16(n uint16) {
	Check(binary.Write(b.writer, binary.BigEndian, n))
}

func (b *Binary) WriteUInt8(n uint8) {
	Check(binary.Write(b.writer, binary.BigEndian, n))
}

func (b *Binary) ReadUInt8() uint8 {
	var ret uint8
	Check(binary.Read(b.reader, binary.BigEndian, &ret))
	return ret
}

func (b *Binary) ReadUInt16() uint16 {
	var ret uint16
	Check(binary.Read(b.reader, binary.BigEndian, &ret))
	return ret
}
