package sockets

type Sockets interface {
	Listen(addr string) string
	Accept(name string) string
	Connect(addr string) string
	Close(id string)
	Read(id string, max int) []byte
	Write(id string, data []byte)
}
