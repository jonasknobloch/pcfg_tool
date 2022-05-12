package pcfg

import (
	"errors"
	"fmt"
	"github.com/jonasknobloch/jinn/pkg/tree"
)

type Span struct {
	i, j int
	n    string
}

type Item struct {
	Span
	weight     float64
	backtracks [2]*Item
}

func (i *Item) Weight() float64 {
	return i.weight
}

func (i *Item) String() string {
	return fmt.Sprintf("(%d,%s,%d)#%.2f", i.i, i.n, i.j, i.Weight())
}

type Parser struct {
	tokens  []string
	grammar Grammar
	heap    Heap
	matcher Matcher
	rules   map[string][]Rule
}

func NewParser(g *Grammar) (*Parser, error) {
	if g.initial == "" {
		return nil, errors.New("grammar initial not set")
	}

	rules := make(map[string][]Rule)

	add := func(k string, r Rule) {
		if _, ok := rules[k]; !ok {
			rules[k] = make([]Rule, 0)
		}

		rules[k] = append(rules[k], r)
	}

	for r := range g.weights {
		switch v := r.(type) {
		case *Lexical:
			add(v.body, r)
		case *NonLexical:
			add(v.body[0], r)

			if v.body[len(v.body)-1] != v.body[0] {
				add(v.body[len(v.body)-1], r)
			}
		default:
			panic("unknown rule type")
		}
	}

	return &Parser{
		grammar: *g,
		rules:   rules,
	}, nil
}

func (p *Parser) Parse(tokens []string) (*tree.Tree, bool) {
	p.tokens = tokens

	p.heap = *NewHeap()
	p.matcher = *NewMatcher()

	p.Initialize()

	for !p.heap.Empty() {
		item, _ := p.heap.Pop()

		if ok := p.matcher.Add(item); !ok {
			continue
		}

		if item.n == p.grammar.initial && item.i == 0 && item.j == len(p.tokens) {
			return p.Tree(item), true
		}

		rules, ok := p.rules[item.n]

		if !ok {
			continue
		}

		for _, rule := range rules {
			nonLexical, ok := rule.(*NonLexical)

			if !ok {
				continue
			}

			if len(nonLexical.body) == 2 {
				if nonLexical.body[0] == item.n {
					for _, c := range p.matcher.MatchLeft(item.j, nonLexical.body[1]) {
						p.Combine(item, c, rule)
					}
				}

				if nonLexical.body[1] == item.n {
					for _, c := range p.matcher.MatchRight(nonLexical.body[0], item.i) {
						p.Combine(c, item, rule)
					}
				}
			}

			if len(nonLexical.body) == 1 {
				p.Chain(item, rule)
			}
		}
	}

	return nil, false
}

func (p *Parser) Initialize() {
	for i, t := range p.tokens {
		terminal := &Item{
			Span: Span{
				i: i,
				j: i + 1,
				n: t,
			},
			weight: 1,
		}

		for _, r := range p.Rules(t) {
			if _, ok := r.(*Lexical); !ok {
				continue
			}

			lexical := &Item{
				Span: Span{
					i: i,
					j: i + 1,
					n: r.Head(),
				},
				weight:     p.grammar.Weight(r),
				backtracks: [2]*Item{terminal, nil},
			}

			p.heap.Push(lexical)
		}
	}
}

func (p *Parser) Rules(body string) []Rule {
	rules, ok := p.rules[body]

	if !ok {
		rules = []Rule{}
	}

	return rules
}

func (p *Parser) Combine(c1, c2 *Item, r Rule) {
	i := &Item{
		Span: Span{
			i: c1.i,
			j: c2.j,
			n: r.Head(),
		},
		weight:     c1.Weight() * c2.Weight() * p.grammar.Weight(r),
		backtracks: [2]*Item{c1, c2},
	}

	p.heap.Push(i)
}

func (p *Parser) Chain(c *Item, r Rule) {
	i := &Item{
		Span: Span{
			i: c.i,
			j: c.j,
			n: r.Head(),
		},
		weight:     c.Weight() * p.grammar.Weight(r),
		backtracks: [2]*Item{c, nil},
	}

	p.heap.Push(i)
}

func (p *Parser) Tree(root *Item) *tree.Tree {
	var backtrack func(item *Item) *tree.Tree
	backtrack = func(item *Item) *tree.Tree {
		t := &tree.Tree{
			Label:    item.n,
			Children: nil,
		}

		li, ri := item.backtracks[0], item.backtracks[1]

		if li != nil && ri != nil {
			t.Children = []*tree.Tree{
				backtrack(li),
				backtrack(ri),
			}
		}

		if li != nil && ri == nil {
			t.Children = []*tree.Tree{
				backtrack(li),
			}
		}

		return t
	}

	return backtrack(root)
}
