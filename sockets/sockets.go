package sockets

import (
	"log"
)

const CLOSE_NONE = 0
const CLOSE_SCK = 1
const CLOSE_IN = 2
const CLOSE_OUT = 4


type Sockets interface {
	Listen(addr string) string
	Accept(name string) string
	Connect(addr string) string
	Read(id string, max int) []byte
	Write(id string, data []byte, close int)
	Close(id string, resources int)
}



func ReplyServer(scks Sockets, listenId string) {
	log.Printf("[%s] Starting ReplyServer", listenId)
	for true {
		sckid := scks.Accept(listenId)
		if sckid != "" {
			go func() {
				defer scks.Close(sckid, CLOSE_SCK)
				for true {
					data := scks.Read(sckid, 16 * 1204)
					log.Printf("Reply: %v", data)
					scks.Write(sckid, data, CLOSE_NONE)
				}
			}()
		}
	}
}
