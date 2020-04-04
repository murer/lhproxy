package server

import (
	"log"
	"net/http"

	"github.com/murer/lhproxy/sockets"
	"github.com/murer/lhproxy/util"
)

var scks sockets.Sockets

func Start() {
	scks = sockets.GetNative()
	http.HandleFunc("/", Handle)
	log.Printf("Starting server")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	util.Check(err)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Access: %s %s %s", r.RemoteAddr, r.Method, r.URL)
	if r.URL.Method == "GET" && r.URL.Path == "/version.txt" {
		w.Write([]byte(util.Version))
	} else if r.URL.Method == "POST" {
		HandleSockets(w, e)
	} else {
		http.NotFound(w, r)
	}
}

func HandleSockets(w http.ResponseWriter, r *http.Request) {

}
