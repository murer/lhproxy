package server

import (
	// "net/http"
	"log"
	// "github.com/murer/lhproxy/util"
)

type Tunnel struct {
	URL string
	channel chan *Message
}

func NewTunnel(url string) *Tunnel {
	return &Tunnel{
		URL: url,
		channel: make(chan *Message, 2),
	}
}

func (me *Tunnel) Request(req *Message) *Message {
	log.Printf("Sending: %s", req.Name)
	me.channel <- req
	log.Printf("Message sent: %s", req.Name)
	resp := <- me.channel
	log.Printf("Received: %s -> %s", req.Name, resp.Name)
	return resp
}

func (me *Tunnel) Post() {
	// reader := &Reader{}
	// resp, err := http.Post(me.URL, "application/octet-stream", reader)
	// util.Check(err)
	// log.Printf("code: %d", resp.StatusCode)
	log.Printf("redirecting messages")
	// for true {
		req := <- me.channel
		log.Printf("redirect req: %s", req.Name)
		me.channel <- req
	// }
}

type Reader struct {}

func (me *Reader) Read(buf []byte) (int, error) {
	return 0, nil
}
