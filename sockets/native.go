package sockets

type NativeSockets struct {

}

var native = &NativeSockets{

}

func (scks NativeSockets) Listen(addr string) int {
	return 1
}

func GetNative() *NativeSockets {
	return native
}
