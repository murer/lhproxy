package sockets

import (
	"log"
	"net"
	"fmt"
	"strings"

	"github.com/murer/lhproxy/util"
	"github.com/murer/lhproxy/util/queue"
)

type connWrapper struct {
	id string
	conn net.Conn
	reader *queue.Queue
}

func (c *connWrapper) Close() {
	c.conn.Close()
}

func (c *connWrapper) startReading() {
	log.Printf("Starting reading conn: %s", c.id)
}

type listenerWrapper struct {
	id string
	ln net.Listener
	queue *queue.Queue
}

func (l *listenerWrapper) Close() {
	l.ln.Close()
}

func (l *listenerWrapper) accept() *connWrapper {
	ret := l.queue.Shift()
	if ret == nil {
		return nil
	}
	return ret.(*connWrapper)
}

func (l *listenerWrapper) startAccepts() {
	log.Printf("Starting accepts")
	for true {
		conn, err := l.ln.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			util.Check(err)
		}
		c := &connWrapper{
			id: fmt.Sprintf("conn://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String()),
			conn: conn,
			reader: queue.New(100 * 1024),
		}
		log.Printf("Caching accepted conn: %s", c.id)
		c.startReading()
		l.queue.Put(c)
	}
}

var lns = make(map[string]*listenerWrapper)
var conns = make(map[string]*connWrapper)
type NativeSockets struct {}
var native = &NativeSockets{}

func (scks *NativeSockets) Listen(addr string) string {
	ln, err := net.Listen("tcp", addr)
	util.Check(err)
	l := &listenerWrapper{
		ln: ln,
		id: fmt.Sprintf("listen://%s", ln.Addr().String()),
		queue: queue.New(1),
	}
	go l.startAccepts()
	lns[l.id] = l
	log.Printf("Listen %s", l.id)
	log.Printf("[TODO] Close listener: %s", l.id)
	return l.id
}

func (scks *NativeSockets) Accept(name string) string {
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

func (scks *NativeSockets) Connect(addr string) string {
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

func (scks *NativeSockets) Close(id string) {
	l := lns[id]
	if l != nil {
		log.Printf("Closing listen %s", l.id)
		delete(lns, l.id)
		l.Close()
	}
	c := conns[id]
	if c != nil {
		log.Printf("Closing connection %s", c.id)
		delete(conns, c.id)
		c.Close()
	}
}

func (scks *NativeSockets) Read(id string, max int) []byte {
	ret := make([]byte, max)
	c := conns[id]
	n, err := c.conn.Read(ret)
	util.Check(err)
	return ret[:n]
}

func GetNative() *NativeSockets {
	return native
}
