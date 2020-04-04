package server

import (
	"encoding/json"
	"log"
	"github.com/murer/lhproxy/util"
)

type Message struct {
	Name string `json:"name"`
	Headers map[string]string `json:"headers`
	Payload []byte `json:"payload"`
}

func (m *Message) Get(name string) string {
	ret := m.Headers[name]
	if ret == "" {
		log.Panicf("Message header is required: %s", name)
	}
	return ret
}

func MessageEnc(msg *Message) []byte {
	buf, err := json.Marshal(msg)
	util.Check(err)
	return buf
}

func MessageDec(buf []byte) *Message {
	ret := &Message{}
	util.Check(json.Unmarshal(buf, ret))
	return ret
}
