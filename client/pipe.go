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

	scks.Connect(c.Address)

	log.Printf("Start pipe to %s", c.Address)
	_, err := io.Copy(c.Writer, c.Reader)
	util.Check(err)
	log.Printf("Pipe to %s is done", c.Address)
}
