package pcfg

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

	if item.pruned {
		return h.Pop()
	}

	return item, true
}

func (h *Heap) Empty() bool {
	if h.bh.Empty() {
		return true
	}

	return h.bh.Empty()
}

func (h *Heap) Prune(threshold float64) {
	val, ok := h.bh.Peek()

	if !ok {
		return
	}

	if item, ok := val.(*Item); !ok {
		return
	} else {
		threshold *= item.weight
	}

	it := h.bh.Iterator()

	for it.Next() {
		val := it.Value()

		item, ok := val.(*Item)

		if !ok {
			return
		}

		if item.pruned {
			continue
		}

		if item.weight < threshold {
			item.pruned = true
		}
	}
}
