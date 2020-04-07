package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/murer/lhproxy/sockets"
)

func SocksTest(t *testing.T, scks sockets.Sockets) {
	listen := scks.Listen("127.0.0.1:5001")
	assert.NotNil(t, listen)
	defer scks.Close(listen, sockets.CLOSE_SCK)
	assert.Empty(t, scks.Accept(listen))

	cc1 := scks.Connect("127.0.0.1:5001")
	assert.NotEmpty(t, cc1)
	defer scks.Close(cc1, sockets.CLOSE_SCK)

	cc2 := scks.Connect("127.0.0.1:5001")
	assert.NotEmpty(t, cc2)
	defer scks.Close(cc2, sockets.CLOSE_SCK)

	time.Sleep(10 * time.Millisecond)
	cs1 := scks.Accept(listen)
	assert.NotEmpty(t, cs1)
	defer scks.Close(cs1, sockets.CLOSE_SCK)

	time.Sleep(10 * time.Millisecond)
	cs2 := scks.Accept(listen)
	assert.NotEmpty(t, cs2)
	defer scks.Close(cs2, sockets.CLOSE_SCK)

	assert.Equal(t, []byte{}, scks.Read(cs1, 2))
	scks.Write(cc1, []byte{5, 6, 7}, sockets.CLOSE_NONE)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, []byte{5, 6}, scks.Read(cs1, 2))
	assert.Equal(t, []byte{7}, scks.Read(cs1, 2))
	assert.Equal(t, []byte{}, scks.Read(cs1, 2))

	assert.Equal(t, []byte{}, scks.Read(cc1, 2))
	scks.Write(cs1, []byte{5, 6, 7}, sockets.CLOSE_NONE)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, []byte{5, 6}, scks.Read(cc1, 2))
	assert.Equal(t, []byte{7}, scks.Read(cc1, 2))
	assert.Equal(t, []byte{}, scks.Read(cc1, 2))

	scks.Close(cc1, sockets.CLOSE_SCK)
	time.Sleep(100 * time.Millisecond)
	assert.Nil(t, scks.Read(cs1, 2))
	assert.Nil(t, scks.Read(cs1, 2))

	scks.Write(cs2, []byte{1, 2}, sockets.CLOSE_OUT)
	assert.Equal(t, []byte{1, 2}, scks.Read(cc2, 2))
	assert.Nil(t, scks.Read(cc2, 2))
	assert.Equal(t, []byte{}, scks.Read(cs2, 2))
	scks.Close(cc2, sockets.CLOSE_SCK)
	assert.Nil(t, scks.Read(cs2, 2))
}

func GetNative() sockets.Sockets {
	return &sockets.NativeSockets{
		ReadTimeout:       1 * time.Millisecond,
		AcceptTimeout:     1 * time.Millisecond,
	}
}

func TestNativeSockets(t *testing.T) {
	native := GetNative()
	SocksTest(t, native)
}

func TestIdle(t *testing.T) {
	scks := &sockets.NativeSockets{
		ReadTimeout:       1 * time.Millisecond,
		AcceptTimeout:     1 * time.Millisecond,
	}
	go scks.IdleStart(200 * time.Millisecond, 600 * time.Millisecond)

	listen := scks.Listen("127.0.0.1:5001")
	defer scks.Close(listen, sockets.CLOSE_SCK)
	cc := scks.Connect("127.0.0.1:5001")
	defer scks.Close(cc, sockets.CLOSE_SCK)
	cs := scks.Accept(listen)
	defer scks.Close(cs, sockets.CLOSE_SCK)

	scks.Write(cc, []byte{5, 6, 7}, sockets.CLOSE_NONE)
	assert.Equal(t, []byte{5, 6}, scks.Read(cs, 2))

	time.Sleep(400 * time.Millisecond)

	assert.Panics(t, func() {
		assert.Equal(t, []byte{7}, scks.Read(cs, 2))
	})

	assert.Panics(t, func() {
		scks.Read(cc, 2)
	})

	time.Sleep(400 * time.Millisecond)
	assert.Panics(t, func() {
		cs2 := scks.Accept(listen)
		defer scks.Close(cs2, sockets.CLOSE_SCK)
	})
}
