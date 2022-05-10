package pcfg

import (
	"fmt"
	"github.com/jonasknobloch/jinn/pkg/tree"
)

type Span struct {
	i, j int
	n    string
}

type Item struct {
	Span
	p float64
}

func (i *Item) Weight() float64 {
	return i.p
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
	trace   map[*Item][2]*Item
}

func NewParser(g *Grammar) *Parser {
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
	}
}

func (p *Parser) Parse(tokens []string) (*tree.Tree, bool) {
	p.tokens = tokens

	p.heap = *NewHeap()
	p.matcher = *NewMatcher()

	p.trace = make(map[*Item][2]*Item)

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
			lexical, ok := rule.(*NonLexical)

			if !ok {
				continue
			}

			if len(lexical.body) == 2 {
				if lexical.body[0] == item.n {
					for _, c := range p.matcher.MatchLeft(item.j, lexical.body[1]) {
						p.Combine(item, c, rule)
					}
				}

				if lexical.body[1] == item.n {
					for _, c := range p.matcher.MatchRight(lexical.body[0], item.i) {
						p.Combine(c, item, rule)
					}
				}
			}

			if len(lexical.body) == 1 {
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
			p: 1,
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
				p: 1,
			}

			p.heap.Push(lexical)

			p.trace[lexical] = [2]*Item{terminal, nil}
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
		p: c1.Weight() * c2.Weight() * p.grammar.Weight(r),
	}

	p.trace[i] = [2]*Item{c1, c2}

	p.heap.Push(i)
}

func (p *Parser) Chain(c *Item, r Rule) {
	i := &Item{
		Span: Span{
			i: c.i,
			j: c.j,
			n: r.Head(),
		},
		p: c.Weight() * p.grammar.Weight(r),
	}

	p.trace[i] = [2]*Item{c, nil}

	p.heap.Push(i)
}

func (p *Parser) Tree(root *Item) *tree.Tree {
	var trace func(item *Item) *tree.Tree
	trace = func(item *Item) *tree.Tree {
		t := &tree.Tree{
			Label:    item.n,
			Children: nil,
		}

		var li, ri *Item

		if pred, ok := p.trace[item]; !ok {
			return t
		} else {
			li, ri = pred[0], pred[1]
		}

		if ri == nil {
			t.Children = make([]*tree.Tree, 1)
		} else {
			t.Children = make([]*tree.Tree, 2)
		}

		t.Children[0] = trace(li)

		if ri != nil {
			t.Children[1] = trace(ri)
		}

		return t
	}

	return trace(root)
}
