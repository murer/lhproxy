package util

import (
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestQueueInt(t *testing.T) {
	q := New(2)
	q.Put(10)
	q.Put(20)
	assert.Equal(t, 2, len(q.l))
	assert.Equal(t, 10, q.Shift())
	assert.Equal(t, 20, q.Shift())
	assert.Equal(t, 0, len(q.l))
	assert.Nil(t, q.Shift())
}

func TestQueueStruct(t *testing.T) {
	type mystruct struct {
		N string
	}
	q := New(2)
	q.Put(&mystruct{"a"})
	q.Put(&mystruct{"b"})
	assert.Equal(t, 2, len(q.l))
	assert.Equal(t, &mystruct{"a"}, q.Shift())
	assert.Equal(t, &mystruct{"b"}, q.Shift())
	assert.Equal(t, 0, len(q.l))
	assert.Nil(t, q.Shift())
}

func TestQueueAsync(t *testing.T) {
	q := New(2)
	go func() {
		q.Put(10)
		q.Put(20)
		q.Put(30)
	}()
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 10, q.Shift())
	assert.Equal(t, 20, q.Shift())
	assert.Equal(t, 30, q.Shift())
	assert.Nil(t, q.Shift())
}
