package sockets

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/murer/lhproxy/util"
)

const DESC_ERR_NONE = 0
const DESC_ERR_OTHER = 1
const DESC_ERR_EOF = 2
const DESC_ERR_TIMEOUT = 3
const DESC_ERR_CLOSED = 4

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
	if strings.Contains(err.Error(), "use of closed network connection") {
		return DESC_ERR_CLOSED
	}
	return DESC_ERR_OTHER
}

type connWrapper struct {
	id       string
	conn     *net.TCPConn
	lastUsed time.Time
}

func (c *connWrapper) Close() {
	c.conn.Close()
}

func (c *connWrapper) GetLastUsed() time.Time {
	return c.lastUsed
}

type listenerWrapper struct {
	id       string
	ln       net.Listener
	lastUsed time.Time
}

func (c *listenerWrapper) GetLastUsed() time.Time {
	return c.lastUsed
}

func (l *listenerWrapper) Close() {
	l.ln.Close()
}

func (l *listenerWrapper) accept(timeout time.Duration) (*connWrapper, bool) {
	l.ln.(*net.TCPListener).SetDeadline(time.Now().Add(timeout))
	conn, err := l.ln.Accept()
	derr := DescError(err)
	if derr == DESC_ERR_TIMEOUT {
		return nil, false
	}
	if derr == DESC_ERR_CLOSED {
		return nil, true
	}
	util.Check(err)
	c := &connWrapper{
		id:       fmt.Sprintf("conn://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String()),
		conn:     conn.(*net.TCPConn),
		lastUsed: time.Now(),
	}
	return c, false
}

var lns = make(map[string]*listenerWrapper)
var conns = make(map[string]*connWrapper)

type NativeSockets struct {
	ReadTimeout time.Duration
	AcceptTimeout time.Duration
	SocketIdleTimeout time.Duration
	AcceptIdleTimeout time.Duration
}

func (scks *NativeSockets) Listen(addr string) string {
	ln, err := net.Listen("tcp", addr)
	util.Check(err)
	l := &listenerWrapper{
		ln:       ln,
		id:       fmt.Sprintf("listen://%s", ln.Addr().String()),
		lastUsed: time.Now(),
	}
	lns[l.id] = l
	go scks.idleControl(l.id)
	log.Printf("[%s] Listening", l.id)
	return l.id
}

func (scks *NativeSockets) Accept(name string) string {
	l := lns[name]
	l.lastUsed = time.Now()
	log.Printf("[%s] Accepting", l.id)
	conn, closed := l.accept(scks.AcceptTimeout)
	if closed {
		log.Printf("[%s] Socket is closed to accept connection", l.id)
		return "err://closed"
	}
	if conn == nil {
		log.Printf("[%s] No connection accepted", l.id)
		return ""
	}
	conns[conn.id] = conn
	go scks.idleControl(conn.id)
	log.Printf("[%s] Accepted", conn.id)
	return conn.id
}

func (scks *NativeSockets) Connect(addr string) string {
	conn, err := net.Dial("tcp", addr)
	util.Check(err)
	c := &connWrapper{
		id:       fmt.Sprintf("conn://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String()),
		conn:     conn.(*net.TCPConn),
		lastUsed: time.Now(),
	}
	conns[c.id] = c
	go scks.idleControl(c.id)
	log.Printf("[%s] Connected", c.id)
	return c.id
}

func (scks *NativeSockets) Close(id string, resources int) {
	l := lns[id]
	if l != nil {
		if resources == CLOSE_SCK {
			log.Printf("[%s] Closing listen", l.id)
			delete(lns, l.id)
			l.Close()
		}
	}
	c := conns[id]
	if c != nil {
		if resources == CLOSE_IN {
			log.Printf("[%s] Closing conn reader", c.id)
			c.conn.CloseRead()
		}
		if resources == CLOSE_OUT {
			log.Printf("[%s] Closing conn writer", c.id)
			c.conn.CloseWrite()
		}
		if resources == CLOSE_SCK {
			log.Printf("[%s] Closing conn socket", c.id)
			delete(conns, c.id)
			c.Close()
		}
	}
}

func (scks *NativeSockets) Read(id string, max int) []byte {
	c := conns[id]
	buf := make([]byte, max)
	c.conn.SetReadDeadline(time.Now().Add(scks.ReadTimeout))
	n, err := c.conn.Read(buf)
	derr := DescError(err)
	if derr == DESC_ERR_EOF {
		log.Printf("[%s] Socket EOF...", c.id)
		return nil
	}
	if derr != DESC_ERR_TIMEOUT {
		util.Check(err)
	}
	buf = buf[:n]
	c.lastUsed = time.Now()
	log.Printf("[%s] Read: %d", c.id, len(buf))
	return buf
}

func (scks *NativeSockets) Write(id string, data []byte, close int) {
	c := conns[id]
	c.lastUsed = time.Now()
	log.Printf("[]%s] Write: %d", c.id, len(data))
	if len(data) > 0 {
		n, err := c.conn.Write(data)
		util.Check(err)
		if n != len(data) {
			log.Panicf("[%s] Wrong: %d, should was: %d", c.id, n, len(data))
		}
	}
	scks.Close(id, close)
}
