package grammar

import (
	"errors"
	"github.com/jonasknobloch/jinn/pkg/tree"
)

type Rule interface {
	KeyHead() string
	KeyBody() string
	String() string
}

func NewRule(t *tree.Tree) (Rule, error) {
	if len(t.Children) == 0 {
		return nil, errors.New("tree has no children")
	}

	if len(t.Children[0].Children) == 0 {
		return &Lexical{
			Head: t.Label,
			Body: t.Children[0].Label,
		}, nil
	}

	body := func(ts []*tree.Tree) []string {
		ls := make([]string, len(ts))

		for i, t := range ts {
			ls[i] = t.Label
		}

		return ls
	}(t.Children)

	return &NonLexical{
		Head: t.Label,
		Body: body,
	}, nil
}
