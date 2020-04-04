package client

import (
	"testing"
	"time"
	"io"
	// "github.com/stretchr/testify/assert"
	"github.com/murer/lhproxy/sockets"
)

func TestPipe(t *testing.T) {
	scks := &sockets.NativeSockets{
		ReadTimeout: 1 * time.Millisecond,
		AcceptTimeout: 1 * time.Millisecond,
	}
	lr, lw := io.Pipe()
	defer lw.Close()
	defer lr.Close()
	rr, rw := io.Pipe()
	defer rw.Close()
	defer rr.Close()

	sckid := scks.Listen("localhost:5001")
	defer scks.Close(sckid, sockets.CLOSE_SCK)
	p := &Pipe{
		Scks: scks,
		Address: "localhost:5001",
		Reader: lr,
		Writer: rw,
	}
	go p.Execute()
	lw.Write([]byte{1,2}}
}
