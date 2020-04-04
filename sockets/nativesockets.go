package sockets

import (
	"log"
	"net"
	"fmt"
	"time"
	"io"

	"github.com/murer/lhproxy/util"
)

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
	lastUsed int64
}

func (l *listenerWrapper) Close() {
	l.ln.Close()
}

func (l *listenerWrapper) accept(timeout time.Duration) *connWrapper {
	l.ln.(*net.TCPListener).SetDeadline(time.Now().Add(timeout))
	conn, err := l.ln.Accept()
	derr := DescError(err)
	if derr == DESC_ERR_TIMEOUT {
		return nil
	}
	util.Check(err)
	c := &connWrapper{
		id: fmt.Sprintf("conn://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String()),
		conn: conn,
		lastUsed: time.Now().Unix(),
	}
	return c
}

var lns = make(map[string]*listenerWrapper)
var conns = make(map[string]*connWrapper)
type NativeSockets struct {
	ReadTimeout time.Duration
	AcceptTimeout time.Duration
}
var native = &NativeSockets{
	ReadTimeout: 30 * time.Second,
	AcceptTimeout: 30 * time.Second,
}

func (scks *NativeSockets) Listen(addr string) string {
	ln, err := net.Listen("tcp", addr)
	util.Check(err)
	l := &listenerWrapper{
		ln: ln,
		id: fmt.Sprintf("listen://%s", ln.Addr().String()),
		lastUsed: time.Now().Unix(),
	}
	lns[l.id] = l
	log.Printf("[%s] Listening", l.id)
	return l.id
}

func (scks *NativeSockets) Accept(name string) string {
	l := lns[name]
	l.lastUsed = time.Now().Unix()
	log.Printf("[%s] Accepting", l.id)
	conn := l.accept(scks.AcceptTimeout)
	if conn == nil {
		log.Printf("[%s] No connection accepted", l.id)
		return ""
	}
	conns[conn.id] = conn
	log.Printf("[%s] Accepted", conn.id)
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
	log.Printf("[%s] Connected", c.id)
	return c.id
}

func (scks *NativeSockets) Close(id string, resources int) {
	l := lns[id]
	if l != nil {
		log.Printf("[%s] Closing listen", l.id)
		delete(lns, l.id)
		l.Close()
	}
	c := conns[id]
	if c != nil {
		log.Printf("[%s] Closing connection", c.id)
		delete(conns, c.id)
		c.Close()
	}
}

func (scks *NativeSockets) Read(id string, max int) []byte {
	c := conns[id]
	buf := make([]byte, max)
	c.conn.SetReadDeadline(time.Now().Add(scks.ReadTimeout))
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
	log.Printf("[%s] Read: %d", c.id, len(buf))
	return buf
}

func (scks *NativeSockets) Write(id string, data []byte, closeOut bool) {
	c := conns[id]
	c.lastUsed = time.Now().Unix()
	log.Printf("[]%s] Write: %d", c.id, len(data))
	n, err := c.conn.Write(data)
	util.Check(err)
	if n != len(data) {
		log.Panicf("[%s] Wrong: %d, should was: %d", c.id, n, len(data))
	}
}

func GetNative() Sockets {
	return native
}
