package util

import (
	"testing"
  "errors"

  "github.com/stretchr/testify/assert"
)

func TestCheckSuccess(t *testing.T) {
  Check(nil)
  assert.Panics(t, func() { Check(errors.New("mock error")) })
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
