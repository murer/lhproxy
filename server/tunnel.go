package server

import (
	// "net/http"
	"log"
	"sync"
	// "github.com/murer/lhproxy/util"
)

const MSG_MAX = 2

type reply struct {
	req *Message
	resp *Message
}

type Tunnel struct {
	URL string
	channel chan *Message
	mutex *sync.Cond
	msgs []*reply
}

func NewTunnel(url string) *Tunnel {
	return &Tunnel{
		URL: url,
		channel: make(chan *Message, 2),
		mutex: sync.NewCond(&sync.Mutex{}),
		msgs: make([]*reply, 0, MSG_MAX),
	}
}

func (me *Tunnel) Request(req *Message) *Message {
	log.Printf("Sending: %s", req.Name)
	me.mutex.L.Lock()
	defer me.mutex.L.Unlock()
	log.Printf("UUU: %d", len(me.msgs))
	for len(me.msgs) >= MSG_MAX {
		log.Printf("Waiting for slot to request: %d", len(me.msgs))
		me.mutex.Wait()
	}
	rpl:= &reply{req:req}
	me.msgs = append(me.msgs, rpl)
	idx := len(me.msgs)
	for rpl.resp == nil {
		log.Printf("Waiting response: %d", idx)
		me.mutex.Wait()
	}
	return rpl.resp
}

func (me *Tunnel) Post() {
	log.Printf("redirecting messages")
	me.mutex.L.Lock()
	defer me.mutex.L.Unlock()
	for len(me.msgs) <= 0 {
		log.Printf("No messages to post")
		me.mutex.Wait()
	}
	for idx, rpl := range me.msgs {
		log.Printf("Replying %d", idx)
		rpl.resp = rpl.req
	}
	me.msgs = me.msgs[:0]
	me.mutex.Broadcast()
}

func (me *Tunnel) Start() {
	for true {
		me.Post()
	}
}

type Reader struct {}

func (me *Reader) Read(buf []byte) (int, error) {
	return 0, nil
}
