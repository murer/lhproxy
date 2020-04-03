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

func (q *Queue) Put(elements ...interface{}) {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	for len(q.l) >= q.m {
		log.Printf("Producing %v", elements)
		q.c.Wait()
	}
	log.Printf("Produced %v", elements)
	q.l = append(q.l, elements...)
	q.c.Broadcast()
}

func (q *Queue) internalShift(max int) []interface{} {
	if len(q.l) == 0 {
		log.Printf("Nothing to consume")
		return nil
	}
	m := max
	if len(q.l) < max {
		m = len(q.l)
	}
	ret := q.l[0:m]
	q.l = q.l[m:]
	q.c.Broadcast()
	log.Printf("Consumed %v", ret)
	return ret
}

func (q *Queue) Shift() interface{} {
	ret := q.Shiftn(1)
	if len(ret) == 0 {
		return nil
	}
	return ret[0]
}

func (q *Queue) Shiftn(max int) []interface{} {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	return q.internalShift(max)
}

func (q *Queue) WaitShiftn(max int) []interface{} {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	for len(q.l) <= 0 {
		log.Printf("Consuming...")
		q.c.Wait()
	}
	return q.internalShift(1)
}

func (q *Queue) WaitShift() interface{} {
	ret := q.WaitShiftn(1)
	if len(ret) == 0 {
		return nil
	}
	return ret[0]
}

func New(max int) *Queue {
	return &Queue{
		c: sync.NewCond(&sync.Mutex{}),
		m: max,
	}
}
