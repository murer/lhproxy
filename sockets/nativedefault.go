package sockets

import (
	"time"
)

var native = &NativeSockets{
	ReadTimeout: 30 * time.Second,
	AcceptTimeout: 30 * time.Second,
	SocketIdleTimeout: 5 * time.Minute,
	AcceptIdleTimeout: 0 * time.Minute,
}

func GetNative() Sockets {
	return native
}
