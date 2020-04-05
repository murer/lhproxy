package util

import (
	"io"
	"log"
	"os"
	"io/ioutil"
)

var Version = "dev"

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadAll(r io.Reader) []byte {
	ret, err := ioutil.ReadAll(r)
	Check(err)
	return ret
}

func ReadAllString(r io.Reader) string {
	return string(ReadAll(r))
}

func ReadFully(r io.Reader, n int) []byte {
	buf := make([]byte, n)
	nr, err := io.ReadAtLeast(r, buf, n)
	Check(err)
	if nr != n {
		log.Panicf("wrong %d, expected %d", nr ,n)
	}
	return buf
}

func WriteFully(w io.Writer, buf []byte) {
	nr, err := w.Write(buf)
	Check(err)
	if nr != len(buf) {
		log.Panicf("wrong %d, expected %d", nr ,len(buf))
	}
}

func Secret() []byte {
	ret := os.Getenv("LHPROXY_SECRET")
	if ret == "" {
		log.Panicf("LHPROXY_SECRET not found")
	}
	return []byte(ret)
}
