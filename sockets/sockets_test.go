package sockets

import (
	"testing"
	// "net/http"
	// "net/http/httptest"
	//
	"github.com/stretchr/testify/assert"
	//
	// "github.com/murer/lhproxy/util"
)

func TestSockets(t *testing.T) {
	scks := GetNative()
	listen := scks.Listen("127.0.0.1:5001")
	assert.NotNil(t, listen)
	defer scks.Close(listen)
	assert.Empty(t, scks.Accept(listen))
	c := scks.Connect("127.0.0.1:5001")
	assert.NotEmpty(t, c)
	defer scks.Close(c)
}
