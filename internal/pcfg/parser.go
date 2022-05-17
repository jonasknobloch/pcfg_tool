package pcfg

import (
	"errors"
	"fmt"
	"github.com/jonasknobloch/jinn/pkg/tree"
)

type NonLexicalInt struct {
	head int
	body []int
	rule *NonLexical
}

type Span struct {
	i, j int
	n    int
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
	return fmt.Sprintf("(%d,%d,%d)#%.2f", i.i, i.n, i.j, i.Weight())
}

var ErrNoParse = errors.New("no parse")

type Parser struct {
	tokens  []string
	grammar Grammar
	heap    Heap
	matcher Matcher
	symbols SymbolTable
	initial int
	lexicon map[string][]*Lexical
	rules   map[int][]*NonLexicalInt
}

func NewParser(g *Grammar) (*Parser, error) {
	if g.initial == "" {
		return nil, errors.New("grammar initial not set")
	}

	p := &Parser{
		grammar: *g,
		symbols: *NewSymbolTable(),
	}

	p.lexicon = make(map[string][]*Lexical)
	p.rules = make(map[int][]*NonLexicalInt)

	addRule := func(k int, ir *NonLexicalInt) {
		if _, ok := p.rules[k]; !ok {
			p.rules[k] = make([]*NonLexicalInt, 0)
		}

		p.rules[k] = append(p.rules[k], ir)
	}

	if initial, err := p.symbols.Atoi(g.initial); err != nil {
		return nil, err
	} else {
		p.initial = initial
	}

	for r := range g.weights {
		switch v := r.(type) {
		case *Lexical:
			if _, ok := p.lexicon[v.body]; !ok {
				p.lexicon[v.body] = make([]*Lexical, 0)
			}

			p.lexicon[v.body] = append(p.lexicon[v.body], v)
		case *NonLexical:
			ri, err := NonLexicalToNonLexicalInt(v, &p.symbols)

			if err != nil {
				return nil, err
			}

			addRule(ri.body[0], ri)

			if ri.body[len(ri.body)-1] != ri.body[0] {
				addRule(ri.body[len(ri.body)-1], ri)
			}
		default:
			panic("unknown rule type")
		}
	}

	return p, nil
}

func (p *Parser) Parse(tokens []string) (*tree.Tree, error) {
	p.tokens = tokens

	p.heap = *NewHeap()
	p.matcher = *NewMatcher()

	if err := p.Initialize(); err != nil {
		return nil, err
	}

	for !p.heap.Empty() {
		item, _ := p.heap.Pop()

		if ok := p.matcher.Add(item); !ok {
			continue
		}

		if item.n == p.initial && item.i == 0 && item.j == len(p.tokens) {
			return p.Tree(item, tokens)
		}

		rules, ok := p.rules[item.n]

		if !ok {
			continue
		}

		for _, rule := range rules {
			if len(rule.body) == 2 {
				if rule.body[0] == item.n {
					for _, c := range p.matcher.MatchLeft(item.j, rule.body[1]) {
						p.Combine(item, c, rule)
					}
				}

				if rule.body[1] == item.n {
					for _, c := range p.matcher.MatchRight(rule.body[0], item.i) {
						p.Combine(c, item, rule)
					}
				}
			}

			if len(rule.body) == 1 {
				p.Chain(item, rule)
			}
		}
	}

	return nil, ErrNoParse
}

func (p *Parser) Initialize() error {
	for i, t := range p.tokens {
		for _, r := range p.Lexicon(t) {
			n, err := p.symbols.Atoi(r.head)

			if err != nil {
				return err
			}

			lexical := &Item{
				Span: Span{
					i: i,
					j: i + 1,
					n: n,
				},
				weight: p.grammar.Weight(r),
			}

			p.heap.Push(lexical)
		}
	}

	return nil
}

func (p *Parser) Lexicon(body string) []*Lexical {
	lexicon, ok := p.lexicon[body]

	if !ok {
		lexicon = []*Lexical{}
	}

	return lexicon
}

func (p *Parser) Rules(body int) []*NonLexicalInt {
	rules, ok := p.rules[body]

	if !ok {
		rules = []*NonLexicalInt{}
	}

	return rules
}

func (p *Parser) Combine(c1, c2 *Item, ri *NonLexicalInt) {
	i := &Item{
		Span: Span{
			i: c1.i,
			j: c2.j,
			n: ri.head,
		},
		weight:     c1.Weight() * c2.Weight() * p.grammar.Weight(ri.rule),
		backtracks: [2]*Item{c1, c2},
	}

	p.heap.Push(i)
}

func (p *Parser) Chain(c *Item, ri *NonLexicalInt) {
	i := &Item{
		Span: Span{
			i: c.i,
			j: c.j,
			n: ri.head,
		},
		weight:     c.Weight() * p.grammar.Weight(ri.rule),
		backtracks: [2]*Item{c, nil},
	}

	p.heap.Push(i)
}

func (p *Parser) Tree(root *Item, tokens []string) (*tree.Tree, error) {
	var backtrack func(item *Item) (*tree.Tree, error)
	backtrack = func(item *Item) (*tree.Tree, error) {
		t := &tree.Tree{}

		if label, ok := p.symbols.Itoa(item.n); !ok {
			return nil, errors.New("unknown symbol")
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
