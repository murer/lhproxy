package test

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/murer/lhproxy/server"
)

func TestSockets(t *testing.T) {
	server.Config(GetNative(), []byte("12345678901234561234567890123456"))
	svr := httptest.NewServer(http.HandlerFunc(server.Handle))
	defer svr.Close()
	scks := &server.HttpSockets{
		URL: svr.URL,
		Secret: []byte("12345678901234561234567890123456"),
	}
	SocksTest(t, scks)
}
