package server

import (

)

type HttpSockets struct {
	URL string
}

func Send(mreq *Message) *Message {
	breq := MessageEnc(mreq)
	resp, err := http.Post(url, "application/octet-stream", &breq)
	
}

func (scks *HttpSockets) Listen(addr string) string {
	resp := Send(&Message{
			Name: "scks/listen",
			Headers: map[string]string{"addr": addr},
	})
	return resp.Headers["sckid"]
}

func (scks *HttpSockets) Accept(sckid string) string {
	resp := Send(&Message{
			Name: "scks/accepet",
			Headers: map[string]string{"sckid": sckid},
	})
	return resp.Headers["sckid"]
}

func (scks *HttpSockets) Connect(addr string) string {
	resp := Send(&Message{
			Name: "scks/connect",
			Headers: map[string]string{"addr": addr},
	})
	return resp.Headers["sckid"]
}

func (scks *HttpSockets) Close(sckid string, resources int) {
	Send(&Message{
			Name: "scks/close",
			Headers: map[string]string{"sckid": sckid, "crsrc": string(resources)},
	})
}

func (scks *HttpSockets) Read(sckid string, max int) []byte {
	resp := Send(&Message{
			Name: "scks/read",
			Headers: map[string]string{"sckid": sckid,"max":string(max)},
	})
	return resp.Payload
}

func (scks *HttpSockets) Write(sckid string, data []byte, resources int) {
	Send(&Message{
			Name: "scks/write",
			Headers: map[string]string{"sckid": sckid, "crsrc": string(resources)},
			Payload: data,
	})
}
