package client

import (
	"testing"
	"time"
	"io"
	"github.com/stretchr/testify/assert"
	"github.com/murer/lhproxy/sockets"
	"github.com/murer/lhproxy/util"
)

func TestPipe(t *testing.T) {
	scks := &sockets.NativeSockets{
		ReadTimeout: 3000 * time.Millisecond,
		AcceptTimeout: 3000 * time.Millisecond,
	}
	lr, lw := io.Pipe()
	defer lw.Close()
	defer lr.Close()
	rr, rw := io.Pipe()
	defer rw.Close()
	defer rr.Close()

	sckid := scks.Listen("localhost:5001")
	defer scks.Close(sckid, sockets.CLOSE_SCK)
	go sockets.ReplyServer(scks, sckid)

	p := &Pipe{
		Scks: scks,
		Address: "localhost:5001",
		Reader: lr,
		Writer: rw,
	}
	go p.Execute()
	lw.Write([]byte{1,2})
	assert.Equal(t, []byte{1,2}, util.ReadFully(rr, 2))
	lw.Write([]byte{3,4})
	lw.Close()
	assert.Equal(t, []byte{3,4}, util.ReadAll(rr))
}
