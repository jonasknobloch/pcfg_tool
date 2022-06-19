package parser

import (
	"github.com/emirpasic/gods/trees/binaryheap"
)

type Heap struct {
	bh *binaryheap.Heap
}

type HeapItem struct {
	item     *Item
	priority float64
}

func NewHeap() *Heap {
	inverseHeapItemComparator := func(a, b interface{}) int {
		aAsserted := a.(HeapItem)
		bAsserted := b.(HeapItem)

		switch {
		case aAsserted.priority > bAsserted.priority:
			return -1
		case aAsserted.priority < bAsserted.priority:
			return 1
		default:
			return 0
		}
	}

	return &Heap{
		bh: binaryheap.NewWith(inverseHeapItemComparator),
	}
}

func (h *Heap) Push(item *Item, priority float64) {
	h.bh.Push(HeapItem{
		item:     item,
		priority: priority,
	})
}

func (h *Heap) Pop() (*Item, bool) {
	val, ok := h.bh.Pop()

	if !ok {
		return nil, false
	}

	heapItem, ok := val.(HeapItem)

	if !ok {
		panic("unexpected item type")
	}

	return heapItem.item, true
}

func (h *Heap) Empty() bool {
	return h.bh.Empty()
}
