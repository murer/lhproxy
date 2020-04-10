package server

import (

)

type Tunnel struct {}

func (t *Tunnel) Request(r *Message) *Message {
	return r
}

func (t *Tunnel) Post() {

}
