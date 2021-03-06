package grammar

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestSymbolTable(t *testing.T) {
	st := NewSymbolTable()

	for _, c := range []string{"foo", "bar", "baz"} {
		i := st.Atoi(c)

		s, err := st.Itoa(i)

		if err != nil {
			t.Errorf("unexptected error: %v", err)
		}

		if s != c {
			t.Errorf("expected %s got %s", c, s)
		}

		fmt.Println(s, i)
	}
}

func BenchmarkStringComparison(b *testing.B) {
	foo := strconv.Itoa(int(rand.Int63()))
	bar := strconv.Itoa(int(rand.Int63()))

	for i := 0; i < b.N; i++ {
		_ = foo == bar
	}
}

func BenchmarkInt64Comparison(b *testing.B) {
	foo := rand.Int63()
	bar := rand.Int63()

	for i := 0; i < b.N; i++ {
		_ = foo == bar
	}
}
