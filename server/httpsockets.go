package server

import (
	"net/http"
	"bytes"
	"log"
	"io/ioutil"
	"strconv"
	"github.com/murer/lhproxy/util"
)

type HttpSockets struct {
	URL string
	Secret []byte
}

func (scks *HttpSockets) Send(mreq *Message) *Message {
	breq := MessageEnc(scks.Secret, mreq)
	resp, err := http.Post(scks.URL, "application/octet-stream", bytes.NewBuffer(breq))
	util.Check(err)
	if resp.StatusCode != 200 {
		log.Panicf("resp: %s" + resp.Status)
	}
	bresp, err := ioutil.ReadAll(resp.Body)
	util.Check(err)
	mresp := MessageDec(scks.Secret, bresp)
	return mresp
}

func (scks *HttpSockets) Listen(addr string) string {
	resp := scks.Send(&Message{
			Name: "scks/listen",
			Headers: map[string]string{"addr": addr},
	})
	return resp.Headers["sckid"]
}

func (scks *HttpSockets) Accept(sckid string) string {
	resp := scks.Send(&Message{
			Name: "scks/accept",
			Headers: map[string]string{"sckid": sckid},
	})
	return resp.Headers["sckid"]
}

func (scks *HttpSockets) Connect(addr string) string {
	resp := scks.Send(&Message{
			Name: "scks/connect",
			Headers: map[string]string{"addr": addr},
	})
	return resp.Headers["sckid"]
}

func (scks *HttpSockets) Close(sckid string, resources int) {
	scks.Send(&Message{
			Name: "scks/close",
			Headers: map[string]string{"sckid": sckid, "crsrc": strconv.Itoa(resources)},
	})
}

func (scks *HttpSockets) Read(sckid string, max int) []byte {
	resp := scks.Send(&Message{
			Name: "scks/read",
			Headers: map[string]string{"sckid": sckid,"max":strconv.Itoa(max)},
	})
	return resp.Payload
}

func (scks *HttpSockets) Write(sckid string, data []byte, resources int) {
	scks.Send(&Message{
			Name: "scks/write",
			Headers: map[string]string{"sckid": sckid, "crsrc": strconv.Itoa(resources)},
			Payload: data,
	})
}
