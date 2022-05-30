package parser

import "testing"

func TestMatcher_Add(t *testing.T) {
	h := NewMatcher()

	foo := &Item{Span: Span{0, 1, 1}}

	if ok := h.Add(foo); !ok {
		t.Fatalf("expetcted %t but got %t", true, ok)
	}

	if ok := h.Add(foo); ok {
		t.Fatalf("expetcted %t but got %t", false, ok)
	}
}

func TestMatcher_MatchLeft(t *testing.T) {
	bar := &Item{Span: Span{1, 2, 2}}
	baz := &Item{Span: Span{2, 3, 3}}

	h := NewMatcher()

	h.Add(bar)
	h.Add(baz)

	var left []*Item

	left = h.MatchLeft(bar.j, 0)

	if len(left) != 0 {
		t.Fatalf("expected no left matches but got %d", len(left))
	}

	left = h.MatchLeft(bar.j, 3)

	if len(left) != 1 {
		t.Fatalf("expected one left match but got %d", len(left))
	}

	if left[0].Span != baz.Span {
		t.Fatalf("expected %v but got %v", baz.Span, left[0].Span)
	}
}

func TestMatcher_MatchRight(t *testing.T) {
	foo := &Item{Span: Span{0, 1, 1}}
	bar := &Item{Span: Span{1, 2, 2}}

	h := NewMatcher()

	h.Add(foo)
	h.Add(bar)

	var right []*Item

	right = h.MatchRight(1, 0)

	if len(right) != 0 {
		t.Fatalf("expected no right matches but got %d", len(right))
	}

	right = h.MatchRight(1, bar.i)

	if len(right) != 1 {
		t.Fatalf("expected one right match but got %d", len(right))
	}

	if right[0].Span != foo.Span {
		t.Fatalf("expected %v but got %v", foo.Span, right[0].Span)
	}
}
