package pcfg

import (
	"errors"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"strings"
)

type Rule interface {
	Head() string
	Body() string
	String() string
}

type Lexical struct {
	head, body string
}

func (l *Lexical) Head() string {
	return l.head
}

func (l *Lexical) Body() string {
	return l.body
}

func (l *Lexical) String() string {
	return l.Head() + " " + l.Body()
}

type NonLexical struct {
	head string
	body []string
}

func (nl *NonLexical) Head() string {
	return nl.head
}

func (nl *NonLexical) Body() string {
	return strings.Join(nl.body, " ")
}

func (nl *NonLexical) String() string {
	return nl.Head() + " -> " + nl.Body()
}

func NewRule(t *tree.Tree) (Rule, error) {
	if len(t.Children) == 0 {
		return nil, errors.New("tree has no children")
	}

	if len(t.Children[0].Children) == 0 {
		return &Lexical{
			head: t.Label,
			body: t.Children[0].Label,
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
		head: t.Label,
		body: body,
	}, nil
}
