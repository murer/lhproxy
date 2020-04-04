package sockets

import (
	"log"
	"net"
	"fmt"
	"strings"
	"time"
	"io"

	"github.com/murer/lhproxy/util"
	"github.com/murer/lhproxy/util/queue"
)

const READ_DEADLINE = 1 * time.Second

const DESC_ERR_NONE = 0
const DESC_ERR_OTHER = 1
const DESC_ERR_EOF = 2
const DESC_ERR_TIMEOUT = 3

func DescError(err error) int {
	if err == nil {
		return DESC_ERR_NONE
	}
	if err == io.EOF {
		return DESC_ERR_EOF
	}
	netErr, ok := err.(net.Error)
	if ok && netErr.Timeout() {
		return DESC_ERR_TIMEOUT
	}
	return DESC_ERR_OTHER
}

type connWrapper struct {
	id string
	conn net.Conn
	lastUsed int64
}

func (c *connWrapper) Close() {
	c.conn.Close()
}

type listenerWrapper struct {
	id string
	ln net.Listener
	queue *queue.Queue
	lastUsed int64
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
			lastUsed: time.Now().Unix(),
		}
		log.Printf("Caching accepted conn: %s", c.id)
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
		lastUsed: time.Now().Unix(),
	}
	go l.startAccepts()
	lns[l.id] = l
	log.Printf("Listen %s", l.id)
	return l.id
}

func (scks *NativeSockets) Accept(name string) string {
	l := lns[name]
	l.lastUsed = time.Now().Unix()
	log.Printf("Accepting %s", l.id)
	conn := l.accept()
	if conn == nil {
		log.Printf("No connection accepted: %s", l.id)
		return ""
	}
	conns[conn.id] = conn
	log.Printf("Accepted %s", conn.id)
	return conn.id
}

func (scks *NativeSockets) Connect(addr string) string {
	conn, err := net.Dial("tcp", addr)
	util.Check(err)
	c := &connWrapper{
		id: fmt.Sprintf("conn://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String()),
		conn: conn,
		lastUsed: time.Now().Unix(),
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
	c := conns[id]
	buf := make([]byte, max)
	c.conn.SetReadDeadline(time.Now().Add(READ_DEADLINE))
	n, err := c.conn.Read(buf)
	derr := DescError(err)
	if derr == DESC_ERR_EOF {
		log.Printf("Read EOF %s", c.id)
		return nil
	}
	if derr != DESC_ERR_TIMEOUT {
		util.Check(err)
	}
	buf = buf[:n]
	c.lastUsed = time.Now().Unix()
	log.Printf("Read %s: %d", c.id, len(buf))
	return buf
}

func (scks *NativeSockets) Write(id string, data []byte) {
	c := conns[id]
	c.lastUsed = time.Now().Unix()
	log.Printf("Write %s: %d", c.id, len(data))
	n, err := c.conn.Write(data)
	util.Check(err)
	if n != len(data) {
		log.Panicf("Wrong: %d, should was: %d", n, len(data))
	}
}

func GetNative() Sockets {
	return native
}
