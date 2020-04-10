package server

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTunnel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Handle))
	defer server.Close()
	t.Logf("URL: %s", server.URL)
	tunnel := NewTunnel(server.URL)
	defer tunnel.Close()
	original := &Message{Name: "echo", Payload: []byte{10}}
	count := 0
	amount := 1000
	var wg sync.WaitGroup
	for i := 0; i < amount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.Equal(t, original, tunnel.Request(original))
			count++
		}()
	}
	wg.Wait()
	assert.Equal(t, amount, count)
}
