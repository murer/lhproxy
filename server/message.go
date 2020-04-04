package server

type Message struct {
    Name string `json:"name"`
    Headers map[string]string `json:"headers`
    Payload []byte `json:"payload"`
}