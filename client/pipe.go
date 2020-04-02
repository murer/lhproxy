package client

import (
  "io"
	"github.com/murer/lhproxy/util"
)

type Pipe struct {

  RAddress string
  LHAddress string

  LWriter io.Writer
  LReader io.Reader
}

func (c Pipe) Execute() {
  buf := make([]byte, 8 * 1024)
  for true {
    n, err := c.LReader.Read(buf)
    if n > 0 {
      c.LWriter.Write([]byte("LINE: "))
      c.LWriter.Write(buf[:n])
      c.LWriter.Write([]byte{10})
    }
    if err == io.EOF {
      return
    }
    util.Check(err)
  }
}
