package sockets

import (
	"time"
	"reflect"
	"log"
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
	log.Printf("[%s] Starting idle control: %s", name, interval)
	defer log.Printf("[%s] Stoping idle control", name)
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
		log.Printf("[%s] Idle control check lastUsed: %d", name, remaning)
		if remaning <= 0 {
			log.Printf("[%s] Closing idle conn: %d", name, remaning)
			scks.Close(name, CLOSE_SCK)
		}
	}

}
