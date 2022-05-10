package pcfg

import (
	"testing"
)

func TestQueue(t *testing.T) {
	foo := &Item{p: 0.75}
	bar := &Item{p: 0.25}
	baz := &Item{p: 0.5}

	sq := NewQueue()

	sq.Enqueue(foo)
	sq.Enqueue(bar)
	sq.Enqueue(baz)

	for _, g := range []*Item{foo, baz, bar} {
		if i := sq.Dequeue(); i != g {
			t.Errorf("%g expected but got %g", i.Weight(), g.Weight())
		}
	}
}

func BenchmarkReSliceHead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0)

		for j := 0; j < 10000; j++ {
			s = append(s, j)
			s = append(s, j*2)

			s = s[1:]
		}
	}
}

func BenchmarkReSliceTail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0)

		for j := 0; j < 10000; j++ {
			s = append(s, j)
			s = append(s, j*2)

			s = s[:len(s)-1]
		}
	}
}
