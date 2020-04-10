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
  tunnel := &Tunnel{}
	original := &Message{Name: "echo", Payload: []byte{10}}
  assert.Equal(t, original, tunnel.Request(original))
}
