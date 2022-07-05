package parser

import (
	"github.com/emirpasic/gods/trees/redblacktree"
	"unsafe"
)

type RBTree struct {
	t *redblacktree.Tree
}

type RBKey struct {
	item     *Item
	priority float64
}

func NewRBTree() *RBTree {
	rbKeyComparator := func(a, b interface{}) int {
		aAsserted := a.(RBKey)
		bAsserted := b.(RBKey)

		switch {
		case aAsserted.priority > bAsserted.priority:
			return 1
		case aAsserted.priority < bAsserted.priority:
			return -1
		}

		aItemPtr := uintptr(unsafe.Pointer(aAsserted.item))
		bItemPtr := uintptr(unsafe.Pointer(bAsserted.item))

		switch {
		case aItemPtr > bItemPtr:
			return 1
		case aItemPtr < bItemPtr:
			return -1
		}

		return 0
	}

	return &RBTree{t: redblacktree.NewWith(rbKeyComparator)}
}

func (rb *RBTree) Push(item *Item, priority float64) {
	rb.t.Put(RBKey{
		item:     item,
		priority: priority,
	}, nil)
}

func (rb *RBTree) Pop() (*Item, float64, bool) {
	node := rb.t.Right()

	if node == nil {
		return nil, 0, false
	}

	key := node.Key.(RBKey)

	rb.t.Remove(key)

	return key.item, key.priority, true
}

func (rb *RBTree) Empty() bool {
	return rb.t.Empty()
}

func (rb *RBTree) Size() int {
	return rb.t.Size()
}

func (rb *RBTree) Peek() (*Item, float64, bool) {
	node := rb.t.Right()

	if node == nil {
		return nil, 0, false
	}

	key := node.Key.(RBKey)

	return key.item, key.priority, true
}

func (rb *RBTree) Prune(threshold float64) (*Item, bool) {
	node := rb.t.Left()

	if node == nil {
		return nil, false
	}

	key := node.Key.(RBKey)

	if threshold != 0 && key.priority >= threshold {
		return nil, false
	}

	rb.t.Remove(key)

	return key.item, true
}
