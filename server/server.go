package server

import (
	"log"
	"net/http"
	"io/ioutil"
	"github.com/murer/lhproxy/sockets"
	"github.com/murer/lhproxy/util"
)

var scks sockets.Sockets
var secret []byte

func Config(x sockets.Sockets, y []byte) {
	scks = x
	secret = y
}

func Start(addr string, sec []byte) {
	Config(sockets.GetNative(), secret)
	http.HandleFunc("/", Handle)
	log.Printf("Starting server")
	err := http.ListenAndServe(addr, nil)
	util.Check(err)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Access: %s %s %s", r.RemoteAddr, r.Method, r.URL)
	if r.Method == "GET" && r.URL.Path == "/version.txt" {
		w.Write([]byte(util.Version))
	} else if r.Method == "POST" {
		HandleSockets(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func HandleSockets(w http.ResponseWriter, r *http.Request) {
	breq, err := ioutil.ReadAll(r.Body)
	util.Check(err)
	mreq := MessageDec(secret, breq)
	mresp := HandleMessage(mreq)
	bresp := MessageEnc(secret, mresp)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bresp)
}

func HandleMessage(req *Message) *Message {
	if req.Name == "scks/listen" {
		return HandleMessageListen(req)
	} else if req.Name == "scks/accept" {
		return HandleMessageAccept(req)
	} else if req.Name == "scks/connect" {
		return HandleMessageConnect(req)
	} else if req.Name == "scks/write" {
		return HandleMessageWrite(req)
	} else if req.Name == "scks/read" {
		return HandleMessageRead(req)
	} else if req.Name == "scks/close" {
		return HandleMessageClose(req)
	} else {
		log.Panicf("Unknown message %s", req.Name)
	}
	return nil
}

func HandleMessageListen(req *Message) *Message {
	sckid := scks.Listen(req.Get("addr"))
	return &Message{
		Name: "resp/ok",
		Headers: map[string]string{"sckid": sckid},
	}
}

func HandleMessageAccept(req *Message) *Message {
	sckid := scks.Accept(req.Get("sckid"))
	return &Message{
		Name: "resp/ok",
		Headers: map[string]string{"sckid": sckid},
	}
}

func HandleMessageConnect(req *Message) *Message {
	sckid := scks.Connect(req.Get("addr"))
	return &Message{
		Name: "resp/ok",
		Headers: map[string]string{"sckid": sckid},
	}
}

func HandleMessageWrite(req *Message) *Message {
	scks.Write(req.Get("sckid"), req.Payload, req.GetInt("crsrc"))
	return &Message{Name: "resp/ok"}
}

func HandleMessageRead(req *Message) *Message {
	payload := scks.Read(req.Get("sckid"), req.GetInt("max"))
	return &Message{Name: "resp/ok",Payload: payload}
}

func HandleMessageClose(req *Message) *Message {
	scks.Close(req.Get("sckid"), req.GetInt("crsrc"))
	return &Message{Name: "resp/ok"}
}
