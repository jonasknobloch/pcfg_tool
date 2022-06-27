package parser

import (
	"testing"
)

func TestRBTree(t *testing.T) {
	rb := NewRBTree()

	foo := &Item{weight: 0.75}
	bar := &Item{weight: 0.25}
	baz := &Item{weight: 0.5}

	rb.Push(foo, foo.weight)
	rb.Push(bar, bar.weight)
	rb.Push(baz, bar.weight)

	for _, g := range []*Item{foo, baz, bar} {
		if i, _ := rb.Pop(); i != g {
			t.Errorf("%g expected but got %g", i.Weight(), g.Weight())
		}
	}
}

func TestRBTree_PushDuplicateItem(t *testing.T) {
	rb := NewRBTree()

	foo := &Item{weight: 0.75}

	rb.Push(foo, foo.weight)
	rb.Push(foo, foo.weight)

	if rb.t.Size() != 1 {
		t.Fatal("unexpected tree size")
	}

	item, ok := rb.Pop()

	if !ok {
		t.Fatal("item not ok")
	}

	if item != foo {
		t.Fatalf("expected %v but got %v", foo, item)
	}
}

func TestRBTree_PushDuplicateWeight(t *testing.T) {
	rb := NewRBTree()

	foo := &Item{weight: 0.75}
	bar := &Item{weight: 0.75}

	rb.Push(foo, foo.weight)
	rb.Push(bar, bar.weight)

	if rb.t.Size() != 2 {
		t.Fatal("unexpected tree size")
	}

	item, ok := rb.Pop()

	if !ok {
		t.Fatal("item not ok")
	}

	if item != foo && item != bar {
		t.Fatalf("expected %v or %v but got %v", foo, bar, item)
	}
}

func TestRBTree_Prune(t *testing.T) {
	rb := NewRBTree()

	foo := &Item{weight: 0.25}
	bar := &Item{weight: 0.5}
	baz := &Item{weight: 0.75}

	rb.Push(foo, foo.weight)
	rb.Push(bar, bar.weight)
	rb.Push(baz, baz.weight)

	if item, ok := rb.Prune(0.25); ok {
		t.Fatalf("expected %v but got %v", nil, item)
	}

	if item, ok := rb.Prune(0.5); !ok || item != foo {
		t.Fatalf("expected %v but got %v", foo, item)
	}
}
