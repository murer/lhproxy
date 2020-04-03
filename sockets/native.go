package sockets

import (
	"log"
	"net"
	"fmt"

	"github.com/murer/lhproxy/util"
)

type listener struct {
	ln net.Listener
	id string
	conn net.Conn
}

func (l listener) nextConn() {
	if l.conn != nil {
		log.Panicf("there already is a connection")
	}
	conn, err := l.ln.Accept()
	util.Check(err)
	ret := fmt.Sprintf("tcp://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String())
	log.Printf("Cached accepted connection %s", ret)
	l.conn = conn
}

func (l listener) accept() net.Conn {
	if l.conn == nil {
		return l.conn
	}
	return nil
}

var lns = make(map[string]*listener)
var conns = make(map[string]net.Conn)

type NativeSockets struct {

}

var native = &NativeSockets{

}

func (scks NativeSockets) Listen(addr string) string {
	ln, err := net.Listen("tcp", addr)
	util.Check(err)
	l := &listener{
		ln: ln,
		id: fmt.Sprintf("listen://%s", ln.Addr().String()),
	}
	lns[l.id] = l
	log.Printf("Listen %s", l.id)
	log.Printf("[TODO] Close listener: %s", l.id)
	return l.id
}

func (scks NativeSockets) Accept(name string) string {
	l := lns[name]
	log.Printf("Accepting %s", l.id)
	conn := l.accept()
	if conn == nil {
		log.Printf("No connection accepted: %s", l.id)
		return ""
	}
	ret := fmt.Sprintf("tcp://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String())
	log.Printf("Accepted %s", ret)
	log.Printf("[TODO] Close accepeted connection: %s", ret)
	return ret
}

func GetNative() *NativeSockets {
	return native
}
