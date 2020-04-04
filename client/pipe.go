package client

import (
	"io"
	"github.com/murer/lhproxy/sockets"
)

type Pipe struct {
	Scks sockets.Sockets
	Address string
	Writer io.Writer
	Reader io.Reader
}

func (c Pipe) Execute() {
	// buf := make([]byte, 8 * 1024)
	// for true {
	// 	n, err := c.LReader.Read(buf)
	// 	if n > 0 {
	// 		c.LWriter.Write([]byte("LINE: "))
	// 		c.LWriter.Write(buf[:n])
	// 		c.LWriter.Write([]byte{10})
	// 	}
	// 	if err == io.EOF {
	// 		return
	// 	}
	// 	util.Check(err)
	// }
}
