package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTunnel(t *testing.T) {
  tunnel := &Tunnel{}
	original := &Message{Name: "echo", Payload: []byte{10}}
  assert.Equal(t, original, tunnel.Request(original))
}
