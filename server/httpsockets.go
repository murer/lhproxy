package server

type HttpSockets struct {

}

func (scks *HttpSockets) Listen(addr string) string {
  return ""
}

func (scks *HttpSockets) Accept(name string) string {
  return ""
}

func (scks *HttpSockets) Connect(addr string) string {
  return ""
}

func (scks *HttpSockets) Close(id string) {
}

func (scks *HttpSockets) Read(id string, max int) []byte {
  return nil
}

func (scks *HttpSockets) Write(id string, data []byte) {
}
