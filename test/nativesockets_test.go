package test

import (
	"testing"
	"time"
	// "net/http"
	// "net/http/httptest"
	//
	"github.com/stretchr/testify/assert"
	//
	"github.com/murer/lhproxy/sockets"
)

func SocksTest(t *testing.T, scks sockets.Sockets) {
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

	assert.Equal(t, []byte{}, scks.Read(cs1, 2))
	scks.Write(cc1, []byte{5, 6, 7})
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, []byte{5, 6}, scks.Read(cs1, 2))
	assert.Equal(t, []byte{7}, scks.Read(cs1, 2))
	assert.Equal(t, []byte{}, scks.Read(cs1, 2))

	assert.Equal(t, []byte{}, scks.Read(cc1, 2))
	scks.Write(cs1, []byte{5, 6, 7})
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, []byte{5, 6}, scks.Read(cc1, 2))
	assert.Equal(t, []byte{7}, scks.Read(cc1, 2))
	assert.Equal(t, []byte{}, scks.Read(cc1, 2))

	scks.Close(cc1)
	time.Sleep(100 * time.Millisecond)
	assert.Nil(t, scks.Read(cs1, 2))
	assert.Nil(t, scks.Read(cs1, 2))
}

func TestNativeSockets(t *testing.T) {
	SocksTest(t, sockets.GetNative())
}
