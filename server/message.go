package server

import (
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
	if msg.Headers == nil {
		msg.Headers = map[string]string{}
	}
	if msg.Payload == nil {
		msg.Payload = []byte{}
	}
	b := util.NewBinary([]byte{})
	b.WriteString(msg.Name)
	b.WriteUInt16(uint16(len(msg.Headers)))
	for key, value := range msg.Headers {
		b.WriteString(key)
		b.WriteString(value)
	}
	b.WriteBytes(msg.Payload)
	ret := b.Bytes()
	// log.Printf("WRITE %x", ret)
	return ret
}

func rawMessageDec(buf []byte) *Message {
	ret := &Message{Headers:map[string]string{}}
	b := util.NewBinary(buf)
	ret.Name = b.ReadString()
	mapLen := int(b.ReadUInt16())
	for i := 0; i < mapLen; i++ {
		key := b.ReadString()
		value := b.ReadString()
		ret.Headers[key] = value
	}
	ret.Payload = b.ReadBytes()
	// log.Printf("READ %v", ret)
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
