package util

import (
	"io"
	"bytes"
	"encoding/binary"
)

func BinEnc(data interface{}) []byte {
	buf := new(bytes.Buffer)
	BinWrite(buf, data)
	return buf.Bytes()
}

func BinDec(data []byte, ret interface{}) {
	buf := bytes.NewReader(data)
	BinRead(buf, ret)
}

func BinWrite(w io.Writer, data interface{}) {
	Check(binary.Write(w, binary.BigEndian, data))
}

func BinRead(r io.Reader, ret interface{}) {
	Check(binary.Read(r, binary.BigEndian, ret))
}
