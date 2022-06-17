package parser

import "testing"

func TestItemMatcher_Add(t *testing.T) {
	m := NewItemMatcher()

	foo := &Item{Span: Span{0, 1, 1}}

	if ok := m.Add(foo); !ok {
		t.Fatalf("expetcted %t but got %t", true, ok)
	}

	if ok := m.Add(foo); ok {
		t.Fatalf("expetcted %t but got %t", false, ok)
	}
}

func TestItemMatcher_Match(t *testing.T) {
	m := NewItemMatcher()

	foo := &Item{Span: Span{0, 1, 1}}
	bar := &Item{Span: Span{1, 2, 2}}
	baz := &Item{Span: Span{2, 3, 3}}

	m.Add(foo)
	m.Add(bar)
	m.Add(baz)

	left, right := m.Match(bar)

	if ll, lr := len(left), len(right); ll != 1 || lr != 1 {
		t.Fatalf("expexted [%d][%d] but got [%d][%d]", 1, 1, ll, lr)
	}

	if li, ri := left[0], right[0]; li != foo || ri != baz {
		t.Fatalf("expexted %v %v but got %v %v", foo, baz, li, ri)
	}
}
