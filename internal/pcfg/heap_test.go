package pcfg

import (
	"testing"
)

func TestHeap(t *testing.T) {
	foo := &Item{p: 0.75}
	bar := &Item{p: 0.25}
	baz := &Item{p: 0.5}

	h := NewHeap()

	h.Push(foo)
	h.Push(bar)
	h.Push(baz)

	for _, g := range []*Item{foo, baz, bar} {
		if i, _ := h.Pop(); i != g {
			t.Errorf("%g expected but got %g", i.Weight(), g.Weight())
		}
	}
}
