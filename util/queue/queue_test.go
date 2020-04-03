package util

import (
	// "time"
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


	// go func() {
	// 	for i := 1; i < 3; i++ {
	// 		time.Sleep(10 * time.Millisecond)
	// 		g.Put(&i)
	// 	}
	// }()
	//
	// assert.Nil(t, g.Shift())
	// time.Sleep(50 * time.Millisecond)
	// assert.Equal(t, 1, g.Shift())
	// assert.Equal(t, 2, g.Shift())
	// assert.Nil(t, g.Shift())
	//
	// time.Sleep(50 * time.Millisecond)
	// assert.Equal(t, 3, g.Shift())
	// assert.Nil(t, g.Shift())
}
