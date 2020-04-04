package client

import (
	"io"
	"log"
	"github.com/murer/lhproxy/sockets"
	"github.com/murer/lhproxy/util"
)

type Pipe struct {
	Scks sockets.Sockets
	Address string
	Writer io.WriteCloser
	Reader io.ReadCloser
}

func (c Pipe) Execute() {
	defer c.Writer.Close()
	log.Printf("Start pipe to %s", c.Address)
	_, err := io.Copy(c.Writer, c.Reader)
	util.Check(err)
	log.Printf("Pipe to %s is done", c.Address)
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
