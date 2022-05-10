package pcfg

import "sort"

type Weighted interface {
	Weight() float64
}

type Queue []Weighted

func NewQueue() *Queue {
	q := make(Queue, 0)

	return &q
}

func (q *Queue) Len() int {
	return len(*q)
}

func (q *Queue) Less(i, j int) bool {
	return (*q)[i].Weight() < (*q)[j].Weight()
}

func (q *Queue) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

func (q *Queue) Enqueue(i Weighted) {
	*q = append(*q, i)
}

func (q *Queue) Dequeue() Weighted {
	sort.Sort(q)

	i := (*q)[len(*q)-1]
	*q = (*q)[:len(*q)-1]

	return i
}
