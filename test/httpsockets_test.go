package test

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/murer/lhproxy/server"
)

func TestSockets(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(server.Handle))
	defer svr.Close()
	scks := &server.HttpSockets{

	}
	SocksTest(t, scks)
}
