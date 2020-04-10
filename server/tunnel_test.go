package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestTunnel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Handle))
	defer server.Close()
	t.Logf("URL: %s", server.URL)
	tunnel := NewTunnel(server.URL)
	original := &Message{Name: "echo", Payload: []byte{10}}
	go func() {
		time.Sleep(1 * time.Millisecond)
		tunnel.Post()
	}()
	assert.Equal(t, original, tunnel.Request(original))
}
