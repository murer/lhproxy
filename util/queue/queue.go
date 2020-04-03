package util

import (
	"log"
	"sync"
)

type queue struct {
	l []interface{}
	m sync.Mutex
}

func (q *queue) Put(element interface{}) {
	q.m.Lock()
	defer q.m.Unlock()
	log.Printf("Produce %s", element)
	q.l = append(q.l, element)
}

func (q *queue) Shift() interface{} {
	q.m.Lock()
	defer q.m.Unlock()
	if len(q.l) == 0 {
		log.Printf("Nothing to consume")
		return nil
	}
	ret := q.l[0]
	q.l = q.l[1:]
	log.Printf("Consuming %s", ret)
	return ret
}

func New(max int) *queue {
	return &queue{}
}
