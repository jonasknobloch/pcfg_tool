package parser

type Queue interface {
	Push(*Item, float64)
	Pop() (*Item, bool)
	Empty() bool
}
