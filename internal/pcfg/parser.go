package pcfg

import (
	"fmt"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"strings"
)

type Item struct {
	i, j int
	n    string
	p    float64
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
	queue   Queue
	matcher Matcher
	rules   map[string][]Rule
	trace   map[*Item][2]*Item
}

func NewParser(g *Grammar) *Parser {
	rules := make(map[string][]Rule)

	for r := range g.weights {
		b := r.Body()

		if _, ok := rules[b]; !ok {
			rules[b] = make([]Rule, 0)
		}

		rules[b] = append(rules[b], r)
	}

	return &Parser{
		grammar: *g,
		rules:   rules,
	}
}

func (p *Parser) Parse(tokens []string) (*tree.Tree, bool) {
	p.tokens = tokens

	p.queue = *NewQueue()
	p.matcher = *NewMatcher()

	p.trace = make(map[*Item][2]*Item)

	p.Initialize()

	for len(p.queue) > 0 {
		item := p.queue.Dequeue().(*Item)

		if ok := p.matcher.Add(item); !ok {
			continue
		}

		if item.n == p.grammar.initial && item.i == 0 && item.j == len(p.tokens) {
			return p.Tree(item), true
		}

		left, right := p.matcher.Match(item)

		for _, li := range left {
			for _, r := range p.Rules(li.n, item.n) {
				p.Combine(li, item, r)
			}
		}

		for _, ri := range right {
			for _, r := range p.Rules(item.n, ri.n) {
				p.Combine(item, ri, r)
			}
		}

		for _, r := range p.Rules(item.n) {
			p.Chain(item, r)
		}
	}

	return nil, false
}

func (p *Parser) Initialize() {
	for i, t := range p.tokens {
		terminal := &Item{
			i: i,
			j: i + 1,
			n: t,
			p: 1,
		}

		for _, r := range p.Rules(t) {
			lexical := &Item{
				i: i,
				j: i + 1,
				n: r.Head(),
				p: 1,
			}

			p.queue.Enqueue(lexical)

			p.trace[lexical] = [2]*Item{terminal, nil}
		}
	}
}

func (p *Parser) Rules(args ...string) []Rule {
	var body string

	switch len(args) {
	case 0:
		panic("args should not be empty")
	case 1:
		body = args[0]
	default:
		body = strings.Join(args, " ")
	}

	rules, ok := p.rules[body]

	if !ok {
		rules = []Rule{}
	}

	return rules
}

func (p *Parser) Combine(c1, c2 *Item, r Rule) {
	i := &Item{
		i: c1.i,
		j: c2.j,
		n: r.Head(),
		p: c1.Weight() * c2.Weight() * p.grammar.Weight(r),
	}

	p.trace[i] = [2]*Item{c1, c2}

	p.queue.Enqueue(i)
}

func (p *Parser) Chain(c *Item, r Rule) {
	i := &Item{
		i: c.i,
		j: c.j,
		n: r.Head(),
		p: c.Weight() * p.grammar.Weight(r),
	}

	p.trace[i] = [2]*Item{c, nil}

	p.queue.Enqueue(i)
}

func (p *Parser) Tree(root *Item) *tree.Tree {
	var trace func(*Item) *tree.Tree
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
