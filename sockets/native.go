package sockets

import (
	"log"
	"net"
	"fmt"
	"sync"
	"github.com/murer/lhproxy/util"
)

type connWrapper struct {
	id string
	conn net.Conn
}

type listenerWrapper struct {
	id string
	ln net.Listener
	conn *connWrapper
	mutex sync.Mutex
}

func (l listenerWrapper) accept() *connWrapper {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.conn == nil {
		return nil
	}
	ret := l.conn
	go l.nextAccpet()
	return ret
}

func (l listenerWrapper) nextAccpet() {
	conn, err := l.ln.Accept()
	util.Check(err)
	c := &connWrapper{
		id: fmt.Sprintf("conn://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String()),
		conn: conn,
	}
	log.Printf("Caching accepted conn: %s", c.id)
	l.conn = c
}

var lns = make(map[string]*listenerWrapper)
var conns = make(map[string]*connWrapper)

type NativeSockets struct {

}

var native = &NativeSockets{

}

func (scks NativeSockets) Listen(addr string) string {
	ln, err := net.Listen("tcp", addr)
	util.Check(err)
	l := &listenerWrapper{
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
	conns[conn.id] = conn
	log.Printf("Accepted %s", conn.id)
	log.Printf("[TODO] Close accepeted connection: %s", conn.id)
	return conn.id
}

func (scks NativeSockets) Connect(addr string) string {
	conn, err := net.Dial("tcp", addr)
	util.Check(err)
	c := &connWrapper{
		id: fmt.Sprintf("conn://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String()),
		conn: conn,
	}
	conns[c.id] = c
	log.Printf("Connected: %s", c.id)
	return c.id
}

func GetNative() *NativeSockets {
	return native
}
