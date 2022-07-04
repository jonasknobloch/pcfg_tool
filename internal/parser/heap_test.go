package parser

import (
	"testing"
)

func TestHeap(t *testing.T) {
	foo := &Item{weight: 0.75}
	bar := &Item{weight: 0.25}
	baz := &Item{weight: 0.5}

	h := NewHeap()

	h.Push(foo, foo.weight)
	h.Push(bar, bar.weight)
	h.Push(baz, baz.weight)

	for _, g := range []*Item{foo, baz, bar} {
		if i, _, _ := h.Pop(); i != g {
			t.Errorf("%g expected but got %g", i.Weight(), g.Weight())
		}
	}
}
