package parser

import (
	"github.com/emirpasic/gods/trees/binaryheap"
)

type Heap struct {
	bh *binaryheap.Heap
}

func NewHeap() *Heap {
	inverseItemComparator := func(a, b interface{}) int {
		aAsserted := a.(*Item)
		bAsserted := b.(*Item)

		switch {
		case aAsserted.Weight() > bAsserted.Weight():
			return -1
		case aAsserted.Weight() < bAsserted.Weight():
			return 1
		default:
			return 0
		}
	}

	return &Heap{
		bh: binaryheap.NewWith(inverseItemComparator),
	}
}

func (h *Heap) Push(item *Item) {
	h.bh.Push(item)
}

func (h *Heap) Pop() (*Item, bool) {
	val, ok := h.bh.Pop()

	if !ok {
		return nil, false
	}

	item, ok := val.(*Item)

	if !ok {
		panic("unexpected item type")
	}

	return item, true
}

func (h *Heap) Empty() bool {
	return h.bh.Empty()
}
