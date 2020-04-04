package server

import (
	"encoding/json"
	"log"
	"strconv"
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

func (m *Message) GetInt(name string) int {
	ret, err := strconv.Atoi(m.Get(name))
	util.Check(err)
	return ret
}

func rawMessageEnc(msg *Message) []byte {
	buf, err := json.Marshal(msg)
	util.Check(err)
	return buf
}

func rawMessageDec(buf []byte) *Message {
	ret := &Message{}
	util.Check(json.Unmarshal(buf, ret))
	return ret
}

func MessageEnc(msg *Message) []byte {
	raw := rawMessageEnc(msg)
	cryptor := &util.Cryptor{Secret:[]byte("12345678901234561234567890123456")}
	return cryptor.Encrypt(raw)
}

func MessageDec(buf []byte) *Message {
	cryptor := &util.Cryptor{Secret:[]byte("12345678901234561234567890123456")}
	raw := cryptor.Decrypt(buf)
	return rawMessageDec(raw)
}
