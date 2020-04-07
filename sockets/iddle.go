package sockets

import (
	"log"
	// "reflect"
	"time"
)

func computeIdleInterval(scks *NativeSockets) time.Duration {
	interval := scks.SocketIdleTimeout
	if interval == 0 {
		interval = scks.ListenIdleTimeout
	} else if scks.ListenIdleTimeout > 0 && scks.ListenIdleTimeout < interval {
		interval = scks.ListenIdleTimeout
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

func (scks *NativeSockets) idleCheck() {
	for _, conn := range conns {
		scks.closeIfIdle(conn.id, conn.lastUsed, scks.SocketIdleTimeout)
	}
	for _, conn := range lns {
		scks.closeIfIdle(conn.id, conn.lastUsed, scks.ListenIdleTimeout)
	}
}

func (scks *NativeSockets) IdleStart() {
	interval := computeIdleInterval(scks)
	if interval == 0 {
		return
	}
	for true {
		time.Sleep(interval)
		scks.idleCheck()
	}
}
