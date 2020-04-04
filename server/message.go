package server

import (
	"encoding/json"

	"github.com/murer/lhproxy/util"
)

type Message struct {
	Name string `json:"name"`
	Headers map[string]string `json:"headers`
	Payload []byte `json:"payload"`
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
