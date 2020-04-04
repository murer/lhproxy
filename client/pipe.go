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

	sckid string
}

func (c Pipe) local2Server() {
	for true {
		buf := make([]byte, 16 * 1024)
		n, err := c.Reader.Read(buf)
		util.Check(err)
		if n > 0 {
			c.Scks.Write(c.sckid, buf[:n], sockets.CLOSE_NONE)
		}
	}
}

func (c Pipe) server2Local() {
	for true {
		buf := c.Scks.Read(c.sckid, 16 * 1024)
		if buf == nil {
			c.Writer.Close()
			return
		}
		util.WriteFully(c.Writer, buf)
	}
}

func (c Pipe) Execute() {
	defer c.Writer.Close()
	c.sckid = c.Scks.Connect(c.Address)
	go c.local2Server()
	go c.server2Local()
	log.Printf("Pipe to %s is done", c.Address)

	// log.Printf("Start pipe to %s", c.Address)
	// _, err := io.Copy(c.Writer, c.Reader)
	// util.Check(err)
	// log.Printf("Pipe to %s is done", c.Address)
}
