package server

import (
	"log"
	"sync"
)

const MSG_MAX = 20

type reply struct {
	req  *Message
	resp *Message
}

type Tunnel struct {
	URL     string
	channel chan *Message
	mutex   *sync.Cond
	msgs    []*reply
	closed  bool
}

func NewTunnel(url string) *Tunnel {
	ret := &Tunnel{
		URL:     url,
		channel: make(chan *Message, 2),
		mutex:   sync.NewCond(&sync.Mutex{}),
		msgs:    make([]*reply, 0, MSG_MAX),
		closed:  false,
	}
	go ret.start()
	return ret
}

func (me *Tunnel) Request(req *Message) *Message {
	if req == nil {
		log.Panicf("You can not request nil message")
	}
	me.mutex.L.Lock()
	defer me.mutex.L.Unlock()
	if me.closed {
		log.Panicf("Tunnel is closed")
	}
	for len(me.msgs) >= MSG_MAX {
		me.mutex.Wait()
	}
	rpl := &reply{req: req}
	me.msgs = append(me.msgs, rpl)
	for rpl.resp == nil {
		me.mutex.Broadcast()
		me.mutex.Wait()
	}
	return rpl.resp
}

func (me *Tunnel) send() {
	me.mutex.L.Lock()
	defer me.mutex.L.Unlock()
	for len(me.msgs) <= 0 {
		me.mutex.Wait()
	}
	// log.Printf("Posting messages: %d", len(me.msgs))
	if me.msgs[0] == nil {
		log.Printf("Message nil found, stopping...")
		me.mutex.Broadcast()
		me.closed = true
		return
	}
	me.post()
	for _, rpl := range me.msgs {
		if rpl == nil {
			break
		}
		rpl.resp = rpl.req
	}
	me.msgs = me.msgs[:0]
	me.mutex.Broadcast()
}

func (me *Tunnel) start() {
	for !me.closed {
		me.send()
	}
	log.Printf("Posts stopped")
}

func (me *Tunnel) Close() error {
	log.Printf("Closing tunnel")
	me.mutex.L.Lock()
	defer me.mutex.L.Unlock()
	me.msgs = me.msgs[:1]
	me.msgs[0] = nil
	for !me.closed {
		log.Printf("Waiting for nil message unlock routines")
		me.mutex.Broadcast()
		me.mutex.Wait()
	}
	log.Printf("Tunnel closed")
	return nil
}

func (me *Tunnel) post() {

}

type Reader struct{}

func (me *Reader) Read(buf []byte) (int, error) {
	return 0, nil
}
