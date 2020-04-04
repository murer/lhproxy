package sockets

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
