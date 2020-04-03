package util

import (
	"log"
	"sync"
)

type queue struct {
	c *sync.Cond
	m int
	l []interface{}
}

func (q *queue) Put(element interface{}) {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	for len(q.l) >= q.m {
		log.Printf("Producing %s", element)
		q.c.Wait()
	}
	log.Printf("Produced %s", element)
	q.l = append(q.l, element)
}

func (q *queue) Shift() interface{} {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	if len(q.l) == 0 {
		log.Printf("Nothing to consume")
		return nil
	}
	ret := q.l[0]
	q.l = q.l[1:]
	q.c.Broadcast()
	log.Printf("Consuming %s", ret)
	return ret
}

func New(max int) *queue {
	return &queue{
		c: sync.NewCond(&sync.Mutex{}),
		m: max,
	}
}
