package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTunnel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Handle))
	defer server.Close()
	t.Logf("URL: %s", server.URL)
	tunnel := NewTunnel(server.URL)
	original := &Message{Name: "echo", Payload: []byte{10}}
	go tunnel.Post()
	assert.Equal(t, original, tunnel.Request(original))
	assert.Equal(t, original, tunnel.Request(original))
}
