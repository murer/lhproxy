package server

import (
	"log"
	"net/http"

	"github.com/murer/lhproxy/util"
)

func Start() {
	http.HandleFunc("/", Handle)
	log.Printf("Starting server")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	util.Check(err)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Access: %s %s %s", r.RemoteAddr, r.Method, r.URL)
	w.Write([]byte(util.Version))
}
