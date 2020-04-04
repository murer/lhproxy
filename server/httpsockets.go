package server

type HttpSockets struct {

}

func Send(msg *Message) *Message {
    return nil
}

func (scks *HttpSockets) Listen(addr string) string {
  resp := Send(&Message{
      Name: "scks/listen", 
      Headers: map[string]string{"addr": addr},
  })
  return resp.Headers["sckid"]
}

func (scks *HttpSockets) Accept(sckid string) string {
  resp := Send(&Message{
      Name: "scks/accepet", 
      Headers: map[string]string{"sckid": sckid},
  })
  return resp.Headers["sckid"]
}

func (scks *HttpSockets) Connect(addr string) string {
  resp := Send(&Message{
      Name: "scks/connect", 
      Headers: map[string]string{"addr": addr},
  })
  return resp.Headers["sckid"]
}

func (scks *HttpSockets) Close(id string) {
  Send(&Message{
      Name: "scks/close", 
  })
}

func (scks *HttpSockets) Read(id string, max int) []byte {
  resp := Send(&Message{
      Name: "scks/read", 
      Headers: map[string]string{"sckid": sckid},
  })
  return resp.Payload
}

func (scks *HttpSockets) Write(sckid string, data []byte) {
  Send(&Message{
      Name: "scks/write", 
      Headers: map[string]string{"sckid": sckid},
      Payload: data,
  })
}
