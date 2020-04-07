package sockets

import (
	"time"
)

var native *NativeSockets = nil

func GetNative() Sockets {
	if native == nil {
		ret := &NativeSockets{
			ReadTimeout:       30 * time.Second,
			AcceptTimeout:     30 * time.Second,
		}
		go ret.IdleStart(5 * time.Minute, 0 * time.Minute)
		native = ret
	}
	return native
}
