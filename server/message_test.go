package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
    msg := &Message{
        Name: "n",
        Headers: map[string]string{"foo": "1", "bar": "2"},
        Payload: []byte{1,2,3},
    }
    assert.Equal(t, "x", msg)
}