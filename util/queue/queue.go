package util

type queue struct {
	l []interface{}
}

func (q *queue) Put(element interface{}) {
	q.l = append(q.l, element)
}

func (q *queue) Shift() interface{} {
	ret := q.l[0]
	q.l = q.l[1:]
	return ret
}

func New(max int) *queue {
	return &queue{}
}
