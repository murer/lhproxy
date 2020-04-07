package sockets

import (
	"time"
)

var native *NativeSockets = nil

func GetNative() Sockets {
	if native == nil {
		log.Printf("Creating defaut NativeSockets")
		native = &NativeSockets{
			ReadTimeout:       30 * time.Second,
			AcceptTimeout:     30 * time.Second,
		}
	}
	return native
}
