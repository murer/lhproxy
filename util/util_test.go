package util

import (
	"testing"
  "errors"
)

func TestCheckSuccess(t *testing.T) {
  Check(nil)
}

func TestCheckError(t *testing.T) {
  defer func() {
    err := recover()
    if err == nil {
      t.Errorf("Panic was required")
    }
  }()
  Check(errors.New("mock error"))
}
