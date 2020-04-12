package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"fmt"
	"testing"
)

func TestTunnel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Handle))
	defer server.Close()
	url := fmt.Sprintf("%s/tunnel", server.URL)
	t.Logf("URL: %s", url)
	tunnel := NewTunnel(url)
	defer tunnel.Close()
	original := &Message{Name: "echo", Payload: []byte{10}}
	count := 0
	amount := 1
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
