package sockets

import (
	"log"
	// "reflect"
	"time"
)

func computeIdleInterval(socketTimeout time.Duration, listenTimeout time.Duration) time.Duration {
	interval := socketTimeout
	if interval == 0 {
		interval = listenTimeout
	} else if listenTimeout > 0 && listenTimeout < interval {
		interval = listenTimeout
	}
	interval = interval / 10
	return interval
}

func (scks *NativeSockets) closeIfIdle(name string, lastModified time.Time, timeout time.Duration) {
	limit := lastModified.Add(timeout)
	remaning := limit.Sub(time.Now())
	if remaning <= 0 {
		log.Printf("[%s] Closing idle conn: %d", name, remaning)
		scks.Close(name, CLOSE_SCK)
	}
}

func (scks *NativeSockets) idleCheck(socketTimeout time.Duration, listenTimeout time.Duration) {
	for _, conn := range scks.conns {
		scks.closeIfIdle(conn.id, conn.lastUsed, socketTimeout)
	}
	for _, conn := range scks.lns {
		scks.closeIfIdle(conn.id, conn.lastUsed, listenTimeout)
	}
}

func (scks *NativeSockets) IdleStart(socketTimeout time.Duration, listenTimeout time.Duration) {
	interval := computeIdleInterval(socketTimeout, listenTimeout)
	if interval == 0 {
		return
	}
	for true {
		time.Sleep(interval)
		scks.idleCheck(socketTimeout, listenTimeout)
	}
}
