package sockets

import (
	"log"
	"net"
	"fmt"

	"github.com/murer/lhproxy/util"
)

type NativeSockets struct {

}

var native = &NativeSockets{

}

var lns = make(map[string]net.Listener)
var conns = make(map[string]net.Conn)

func (scks NativeSockets) Listen(addr string) string {
	ln, err := net.Listen("tcp", addr)
	util.Check(err)
	ret := fmt.Sprintf("listen://%s", ln.Addr().String())
	lns[ret] = ln
	log.Printf("Listen %s", ln.Addr())
	log.Printf("[TODO] Close listener: %s", ret)
	return ret
}

func (scks NativeSockets) Accept(name string) string {
	ln := lns[name]
	log.Printf("Accepting %s", ln.Addr().String())
	conn, err := ln.Accept()
	util.Check(err)
	ret := fmt.Sprintf("tcp://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String())
	log.Printf("Accepted %s", ret)
	return ret
}

func GetNative() *NativeSockets {
	return native
}
