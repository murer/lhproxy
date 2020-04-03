package queue

import (
	"log"
	"sync"
)

type Queue struct {
	c *sync.Cond
	m int
	l []interface{}
}

func (q *Queue) Put(element interface{}) {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	for len(q.l) >= q.m {
		log.Printf("Producing %v", element)
		q.c.Wait()
	}
	log.Printf("Produced %v", element)
	q.l = append(q.l, element)
	q.c.Broadcast()
}

func (q *Queue) internalShift() interface{} {
	if len(q.l) == 0 {
		log.Printf("Nothing to consume")
		return nil
	}
	ret := q.l[0]
	q.l = q.l[1:]
	q.c.Broadcast()
	log.Printf("Consumed %v", ret)
	return ret
}

func (q *Queue) Shift() interface{} {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	return q.internalShift()
}

func (q *Queue) WaitShift() interface{} {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	for len(q.l) <= 0 {
		log.Printf("Consuming...")
		q.c.Wait()
	}
	return q.internalShift()
}

func New(max int) *Queue {
	return &Queue{
		c: sync.NewCond(&sync.Mutex{}),
		m: max,
	}
}
