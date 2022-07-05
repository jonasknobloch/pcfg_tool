package parser

type Queue interface {
	Push(*Item, float64)
	Pop() (*Item, float64, bool)
	Empty() bool
	Size() int
}
