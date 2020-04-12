package server

import (
	"github.com/murer/lhproxy/util"
	"io"
	"log"
	"net/http"
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
	idxnil := -1
	for idx, rpl := range me.msgs {
		if rpl == nil {
			idxnil = idx
		}
	}
	if idxnil == 0 {
		log.Printf("Message nil found, stopping...")
		me.mutex.Broadcast()
		me.closed = true
		return
	}
	if idxnil < 0 {
		idxnil = len(me.msgs)
	}
	if idxnil >= 0 {
		me.post(idxnil)
	}
	// for _, rpl := range me.msgs {
	// 	if rpl == nil {
	// 		break
	// 	}
	// 	rpl.resp = rpl.req
	// }
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

func aaaa(r io.Reader) int {
	buf := util.ReadFully(r, 2)
	b := util.NewBinary(buf)
	return int(b.ReadUInt16())
}

func handleTunnel(w http.ResponseWriter, r *http.Request) {
	data := util.ReadAll(r.Body)
	breq := util.NewBinary(data)
	total := int(breq.ReadUInt16())
	bresp := util.NewBinary([]byte{})
	bresp.WriteUInt16(uint16(total))
	for i := 0; i < total; i++ {
		raw := breq.ReadBytes()
		mreq := rawMessageDec(raw)
		mresp := HandleMessage(mreq)
		raw = rawMessageEnc(mresp)
		bresp.WriteBytes(raw)
	}
	w.Write(bresp.Bytes())
}

func (me *Tunnel) post(max int) {
	pipein, pipeout := io.Pipe()
	go func() {
		defer pipeout.Close()
		b := util.NewBinary([]byte{})
		b.WriteUInt16(uint16(max))
		for _, rpl := range me.msgs[:max] {
			if rpl == nil {
				return
			}
			b.WriteBytes(rawMessageEnc(rpl.req))
			pipeout.Write(b.Bytes())
			b = util.NewBinary([]byte{})
		}
	}()
	resp, err := http.Post(me.URL, "application/octet-stream", pipein)
	util.Check(err)
	data := util.ReadAll(resp.Body)
	b := util.NewBinary(data)
	total := int(b.ReadUInt16())
	for i := 0; i < total; i++ {
		raw := b.ReadBytes()
		msg := rawMessageDec(raw)
		me.msgs[i].resp = msg
	}

}

type Reader struct{}

func (me *Reader) Read(buf []byte) (int, error) {
	return 0, nil
}
