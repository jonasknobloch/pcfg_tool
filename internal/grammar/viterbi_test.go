package grammar

import (
	"math"
	"testing"
)

func TestViterbiScores_Outside(t *testing.T) {
	g := NewGrammar()

	g.SetInitial("ROOT")

	g.AddLexical(NewLexical("A", "a", 1, g.Symbols))
	g.AddLexical(NewLexical("B", "b", 1, g.Symbols))
	g.AddLexical(NewLexical("X", "a", 0.8, g.Symbols))
	g.AddLexical(NewLexical("Y", "b", 0.7, g.Symbols))

	g.AddNonLexical(NewNonLexical("ROOT", []string{"X", "Y"}, 1, g.Symbols))
	g.AddNonLexical(NewNonLexical("X", []string{"A", "X"}, 0.2, g.Symbols))
	g.AddNonLexical(NewNonLexical("Y", []string{"B", "Y"}, 0.3, g.Symbols))

	vs := NewViterbiScores()

	vs.CalcOutside(g)

	outside := map[string]float64{
		"ROOT": 1,
		"A":    0.1119,
		"B":    0.1679,
		"X":    0.7,
		"Y":    0.8,
	}

	for v, w := range vs.outside {
		s, _ := g.Symbols.Itoa(v)

		f := math.Floor(w*1e4) / 1e4
		o := outside[s]

		if f != o {
			t.Errorf("%s: expected %g but got %g\n", s, o, f)
		}
	}
}
