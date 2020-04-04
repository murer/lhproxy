package pipe

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
	channel chan string
}

func (c Pipe) chSend(name string) {
	c.channel  <- name
}

func (c Pipe) local2Server() {
	defer c.chSend("Sender")
	for true {
		buf := make([]byte, 16 * 1024)
		n, err := c.Reader.Read(buf)
		cs := sockets.CLOSE_NONE
		if err == io.EOF {
			cs = sockets.CLOSE_OUT
		} else {
			util.Check(err)
		}
		c.Scks.Write(c.sckid, buf[:n], cs)
		if err == io.EOF {
			return
		}
	}
}

func (c Pipe) server2Local() {
	defer c.chSend("Receiver")
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
	c.channel = make(chan string)
	c.sckid = c.Scks.Connect(c.Address)
	defer c.Scks.Close(c.sckid, sockets.CLOSE_SCK)
	go c.local2Server()
	go c.server2Local()

	done := <- c.channel
	log.Printf("Redirect is done: %s", done)

	done = <- c.channel
	log.Printf("Redirect is done: %s", done)

	log.Printf("Done")

	// log.Printf("Start pipe to %s", c.Address)
	// _, err := io.Copy(c.Writer, c.Reader)
	// util.Check(err)
	// log.Printf("Pipe to %s is done", c.Address)
}
