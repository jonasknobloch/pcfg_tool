package parser

import "fmt"

type Span struct {
	i, j, n int
}

type Item struct {
	Span
	weight     float64
	backtracks [2]*Item
}

func (i *Item) Weight() float64 {
	return i.weight
}

func (i *Item) String() string {
	return fmt.Sprintf("(%d,%d,%d)#%.2f", i.i, i.n, i.j, i.Weight())
}
