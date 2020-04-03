package sockets

import (
	"log"
	"net"
	"fmt"
	"strings"
	"io"

	"github.com/murer/lhproxy/util"
	"github.com/murer/lhproxy/util/queue"
)

const BUFFER_SIZE = 100 * 1024

type connWrapper struct {
	id string
	conn net.Conn
	reader *queue.Queue
	lastUsed int64
	eof bool
}

func (c *connWrapper) Close() {
	c.conn.Close()
}

func (c *connWrapper) startReading() {
	log.Printf("Starting reading conn: %s", c.id)
	buf := make([]byte, 8 * 1024)
	for true {
		n, err := c.conn.Read(buf)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			} else if err == io.EOF {
				c.reader.Put("EOF")
	      return
	    }
			util.Check(err)
		}
		ret := make([]interface{}, n)
		for i := 0; i < n; i++ {
			ret[i] = buf[i]
		}
		log.Printf("Caching read %s: %x", c.id, len(ret))
		if n > 0 {
			c.reader.Put(ret...)
		}
		if n < 0 {
			panic(n)
		}
	}
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
			reader: queue.New(BUFFER_SIZE),
			eof: false,
		}
		log.Printf("Caching accepted conn: %s", c.id)
		go c.startReading()
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
	return conn.id
}

func (scks *NativeSockets) Connect(addr string) string {
	conn, err := net.Dial("tcp", addr)
	util.Check(err)
	c := &connWrapper{
		id: fmt.Sprintf("conn://%s:%s", conn.RemoteAddr().String(), conn.LocalAddr().String()),
		conn: conn,
		reader: queue.New(BUFFER_SIZE),
		eof: false,
	}
	conns[c.id] = c
	go c.startReading()
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
	if c.eof {
		return nil
	}
	buf := c.reader.Shiftn(max)
	ret := make([]byte, len(buf))
	for i := 0; i < len(buf); i++ {
		switch v := buf[i].(type) {
			case string:
				c.eof = true
				continue
			default:
				ret[i] = v.(byte)
		}
	}
	if c.eof {
		ret = ret[:len(ret)-1]
	}
	if c.eof && len(ret) == 0 {
		ret = nil
	}
	log.Printf("Read %s: %d", c.id, len(ret))
	return ret
}

func (scks *NativeSockets) Write(id string, data []byte) {
	c := conns[id]
	log.Printf("Write %s: %d", c.id, len(data))
	n, err := c.conn.Write(data)
	util.Check(err)
	if n != len(data) {
		log.Panicf("Wrong: %d, should was: %d", n, len(data))
	}
}

func GetNative() *NativeSockets {
	return native
}
