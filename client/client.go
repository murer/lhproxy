package client

import (
  "io"
  "os"

	"github.com/murer/lhproxy/util"
)

func Execute(w io.Writer, r io.Reader) {
  buf := make([]byte, 5)
  for true {
    n, err := r.Read(buf)
    if n > 0 {
      w.Write([]byte("LINE: "))
      w.Write(buf[:n])
      w.Write([]{10})
    }
    if err == io.EOF {
      return
    }
    util.Check(err)
  }
}

func ExecuteStd() {
  Execute(os.Stdout, os.Stdin)
}
