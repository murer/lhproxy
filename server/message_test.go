package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
		original := &Message{
				Name: "n",
				Headers: map[string]string{"foo": "1", "bar": "2"},
				Payload: []byte{1,2},
		}
		buf := MessageEnc(original)
		t.Logf("message: %d", len(buf))
		msg := MessageDec(buf)
		assert.Equal(t, original, msg)
		assert.Equal(t, 80, len(buf))
}
