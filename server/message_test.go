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
		secret := []byte("12345678901234561234567890123456")
		buf := MessageEnc(secret, original)
		msg := MessageDec(secret, buf)
		assert.Equal(t, original, msg)
		assert.Equal(t, 48, len(buf))
}
