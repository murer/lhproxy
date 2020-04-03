package sockets

import (
	"testing"
	"time"
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

	cc1 := scks.Connect("127.0.0.1:5001")
	assert.NotEmpty(t, cc1)
	defer scks.Close(cc1)

	cc2 := scks.Connect("127.0.0.1:5001")
	assert.NotEmpty(t, cc2)
	defer scks.Close(cc2)

	time.Sleep(10 * time.Millisecond)
	cs1 := scks.Accept(listen)
	assert.NotEmpty(t, cs1)
	defer scks.Close(cs1)

	time.Sleep(10 * time.Millisecond)
	cs2 := scks.Accept(listen)
	assert.NotEmpty(t, cs2)
	defer scks.Close(cs2)

	// scks.Read(cs, 3)
}
