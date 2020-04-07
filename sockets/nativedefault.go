package sockets

import (
	"time"
)

var native = &NativeSockets{
	ReadTimeout:       30 * time.Second,
	AcceptTimeout:     30 * time.Second,
	SocketIdleTimeout: 5 * time.Minute,
	ListenIdleTimeout: 0 * time.Minute,
}

func GetNative() Sockets {
	return native
}
