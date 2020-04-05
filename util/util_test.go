package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	Check(nil)
	assert.Panics(t, func() { Check(errors.New("mock error")) })
}
