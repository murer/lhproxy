package server

import (
	"log"
	"strconv"
	"github.com/murer/lhproxy/util"
)

type Message struct {
	Name string
	Headers map[string]string
	Payload []byte
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
	b := util.NewBinary([]byte{})
	b.WriteString(msg.Name)
	b.WriteUInt16(uint16(len(msg.Headers)))
	for key, value := range msg.Headers {
		b.WriteString(key)
		b.WriteString(value)
	}
	b.WriteNillableBytes(msg.Payload)
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
	ret.Payload = b.ReadNillableBytes()
	// log.Printf("READ %v", ret)
	return ret
}

func MessageEnc(secret []byte, msg *Message) []byte {
	raw := rawMessageEnc(msg)
	cryptor := &util.Cryptor{Secret:secret}
	return cryptor.Encrypt(raw)
}

func MessageDec(secret []byte, buf []byte) *Message {
	cryptor := &util.Cryptor{Secret:secret}
	raw := cryptor.Decrypt(buf)
	return rawMessageDec(raw)
}
