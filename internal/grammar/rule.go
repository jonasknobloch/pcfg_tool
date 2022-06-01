package grammar

import (
	"errors"
	"github.com/jonasknobloch/jinn/pkg/tree"
)

type Rule interface {
	Weight() float64
	String(st *SymbolTable) (string, error)
}

func NewRule(t *tree.Tree, st *SymbolTable) (Rule, string, error) {
	if len(t.Children) == 0 {
		return nil, "", errors.New("tree has no children")
	}

	if len(t.Children[0].Children) == 0 {
		l, k := NewLexical(t.Label, t.Children[0].Label, 1, st)

		return l, k, nil
	}

	ls := make([]string, len(t.Children))

	for i, st := range t.Children {
		ls[i] = st.Label
	}

	l, k := NewNonLexical(t.Label, ls, 1, st)

	return l, k, nil
}
