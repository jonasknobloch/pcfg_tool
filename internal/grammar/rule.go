package grammar

import (
	"errors"
	"github.com/jonasknobloch/jinn/pkg/tree"
)

type Rule interface {
	Weight() float64
	Key() string
	String(st *SymbolTable) (string, error)
}

func NewRule(t *tree.Tree, st *SymbolTable) (Rule, error) {
	if len(t.Children) == 0 {
		return nil, errors.New("tree has no children")
	}

	if len(t.Children[0].Children) == 0 {
		return NewLexical(t.Label, t.Children[0].Label, 1, st), nil
	}

	ls := make([]string, len(t.Children))

	for i, st := range t.Children {
		ls[i] = st.Label
	}

	return NewNonLexical(t.Label, ls, 1, st), nil
}
