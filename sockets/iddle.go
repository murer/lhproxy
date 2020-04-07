package sockets

import (
	"log"
	"reflect"
	"time"
)

type IdleController interface {
	GetLastUsed() time.Time
}

func computeIdleInterval(scks *NativeSockets) time.Duration {
	interval := scks.SocketIdleTimeout
	if interval == 0 {
		interval = scks.AcceptIdleTimeout
	} else if scks.AcceptIdleTimeout > 0 && scks.AcceptIdleTimeout < interval {
		interval = scks.AcceptIdleTimeout
	}
	interval = interval / 10
	return interval
}

func (scks *NativeSockets) idleControl(name string) {
	interval := computeIdleInterval(scks)
	if interval == 0 {
		return
	}
	for true {
		time.Sleep(interval)
		timeout := scks.SocketIdleTimeout
		var sck IdleController = conns[name]
		if sck == nil || reflect.ValueOf(sck).IsNil() {
			sck = lns[name]
			timeout = scks.AcceptIdleTimeout
		}
		if sck == nil || reflect.ValueOf(sck).IsNil() {
			return
		}
		limit := sck.GetLastUsed().Add(timeout)
		remaning := limit.Sub(time.Now())
		if remaning <= 0 {
			log.Printf("[%s] Closing idle conn: %d", name, remaning)
			scks.Close(name, CLOSE_SCK)
		}
	}

}
