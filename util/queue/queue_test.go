package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQueueInt(t *testing.T) {
	q := New(2)
	q.Put(10, 20)
	assert.Equal(t, 2, len(q.l))
	assert.Equal(t, 10, q.Shift())
	assert.Equal(t, 20, q.Shift())
	assert.Equal(t, 0, len(q.l))
	assert.Nil(t, q.Shift())

	q.Put(3, 4)
	vs := q.Shiftn(3)
	assert.Equal(t, 3, vs[0])
	assert.Equal(t, 4, vs[1])
	assert.Nil(t, q.Shift())
}

func TestQueueStruct(t *testing.T) {
	type mystruct struct {
		N string
	}
	q := New(2)
	q.Put(&mystruct{"a"}, &mystruct{"b"})
	assert.Equal(t, 2, len(q.l))
	assert.Equal(t, &mystruct{"a"}, q.Shift())
	assert.Equal(t, &mystruct{"b"}, q.Shift())
	assert.Equal(t, 0, len(q.l))
	assert.Nil(t, q.Shift())

	q.Put(&mystruct{"c"}, &mystruct{"d"})
	vs := q.Shiftn(3)
	assert.Equal(t, &mystruct{"c"}, vs[0])
	assert.Equal(t, &mystruct{"d"}, vs[1])
	assert.Nil(t, q.Shift())
}

func TestQueueAsync(t *testing.T) {
	q := New(2)
	go func() {
		time.Sleep(30 * time.Millisecond)
		q.Put(10)
		q.Put(20)
		q.Put(30)
	}()
	assert.Equal(t, 10, q.WaitShift())
	assert.Equal(t, 20, q.WaitShift())
	assert.Equal(t, 30, q.WaitShift())
	assert.Nil(t, q.Shift())
}
