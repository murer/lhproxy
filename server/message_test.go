package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
		original := &Message{
				Name: "nm",
				Headers: map[string]string{"foo": "1", "bar": "2"},
				Payload: []byte{1,2},
		}
		buf := MessageEnc(original)
		msg := MessageDec(buf)
		assert.Equal(t, original, msg)
		assert.Equal(t, 48, len(buf))
}
