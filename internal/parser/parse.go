package parser

import (
	"github.com/jonasknobloch/jinn/pkg/tree"
	"pcfg_tool/internal/grammar"
)

type parse struct {
	tokens  []string
	heap    *Heap
	matcher *Matcher
	parser  *Parser
}

func (p *parse) Parse() (*tree.Tree, error) {
	if err := p.Initialize(); err != nil {
		return nil, err
	}

	for !p.heap.Empty() {
		item, _ := p.heap.Pop()

		if ok := p.matcher.Add(item); !ok {
			continue
		}

		if item.n == p.parser.grammar.Initial() && item.i == 0 && item.j == len(p.tokens) {
			return p.Tree(item, p.tokens)
		}

		for _, rule := range p.parser.grammar.Rules(item.n) {
			if len(rule.Body) == 2 {
				if rule.Body[0] == item.n {
					for _, c := range p.matcher.MatchLeft(item.j, rule.Body[1]) {
						p.Combine(item, c, rule)
					}
				}

				if rule.Body[1] == item.n {
					for _, c := range p.matcher.MatchRight(rule.Body[0], item.i) {
						p.Combine(c, item, rule)
					}
				}
			}

			if len(rule.Body) == 1 {
				p.Chain(item, rule)
			}
		}
	}

	return nil, ErrNoParse
}

func (p *parse) Initialize() error {
	for i, t := range p.tokens {
		for _, r := range p.parser.grammar.Lexicon(t) {
			lexical := &Item{
				Span: Span{
					i: i,
					j: i + 1,
					n: r.Head,
				},
				weight: r.Weight(),
			}

			p.heap.Push(lexical)
		}
	}

	return nil
}

func (p *parse) Combine(c1, c2 *Item, ri *grammar.NonLexical) {
	i := &Item{
		Span: Span{
			i: c1.i,
			j: c2.j,
			n: ri.Head,
		},
		weight:     c1.Weight() * c2.Weight() * ri.Weight(),
		backtracks: [2]*Item{c1, c2},
	}

	p.heap.Push(i)
}

func (p *parse) Chain(c *Item, ri *grammar.NonLexical) {
	i := &Item{
		Span: Span{
			i: c.i,
			j: c.j,
			n: ri.Head,
		},
		weight:     c.Weight() * ri.Weight(),
		backtracks: [2]*Item{c, nil},
	}

	p.heap.Push(i)
}

func (p *parse) Tree(root *Item, tokens []string) (*tree.Tree, error) {
	var backtrack func(item *Item) (*tree.Tree, error)
	backtrack = func(item *Item) (*tree.Tree, error) {
		t := &tree.Tree{}

		if label, err := p.parser.grammar.Symbols.Itoa(item.n); err != nil {
			return nil, err
		} else {
			t.Label = label
		}

		li, ri := item.backtracks[0], item.backtracks[1]

		var errL, errR error

		if li != nil && ri != nil {
			t.Children = make([]*tree.Tree, 2)

			t.Children[0], errL = backtrack(li)
			t.Children[1], errR = backtrack(ri)
		}

		if li != nil && ri == nil {
			t.Children = make([]*tree.Tree, 1)

			t.Children[0], errL = backtrack(li)
		}

		if errL != nil {
			return nil, errL
		}

		if errR != nil {
			return nil, errR
		}

		if li == nil && ri == nil {
			t.Children = []*tree.Tree{
				{
					Label: tokens[item.i],
				},
			}
		}

		return t, nil
	}

	return backtrack(root)
}
