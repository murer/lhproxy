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
	port := scks.Listen("s")
	assert.Less(t, 0, port)
}
