package parser

import (
	"errors"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"pcfg_tool/internal/grammar"
)

type parse struct {
	tokens  []string
	queue   Queue
	matcher *Matcher
	grammar *grammar.Grammar
	viterbi *grammar.ViterbiScores
	config  *Config
}

const UnknownToken = "UNK"

var ErrNoRBTree = errors.New("no red-black tree")

func (p *parse) Parse() (*tree.Tree, error) {
	p.Initialize()

	var rb *RBTree

	if t, ok := p.queue.(*RBTree); !ok {
		rb = t
	}

	if p.config.Prune && rb == nil {
		return nil, ErrNoRBTree
	}

	for !p.queue.Empty() {
		item, _, _ := p.queue.Pop()

		if ok := p.matcher.Add(item); !ok {
			continue
		}

		if item.n == p.grammar.Initial() && item.i == 0 && item.j == len(p.tokens) {
			return p.Tree(item, p.tokens)
		}

		threshold := p.config.Threshold

		if p.config.Prune && threshold != 0 {
			_, priority, ok := rb.Peek()

			if ok {
				threshold *= priority
			}
		}

		for _, rule := range p.grammar.Rules(item.n) {
			if len(rule.Body) == 2 {
				if rule.Body[0] == item.n {
					for _, c := range p.matcher.MatchLeft(item.j, rule.Body[1]) {
						p.Push(p.Combine(item, c, rule), threshold)
					}
				}

				if rule.Body[1] == item.n {
					for _, c := range p.matcher.MatchRight(rule.Body[0], item.i) {
						p.Push(p.Combine(c, item, rule), threshold)
					}
				}
			}

			if len(rule.Body) == 1 {
				p.Push(p.Chain(item, rule), threshold)
			}
		}

		if p.config.Prune {
			for p.config.Rank == 0 || rb.Size() > p.config.Rank {
				if _, ok := rb.Prune(threshold); !ok {
					break
				}
			}
		}
	}

	return nil, ErrNoParse
}

func (p *parse) ItemPriority(i *Item) float64 {
	if !p.config.AStar {
		return i.weight
	}

	return i.weight * p.viterbi.Outside(i.n)
}

func (p *parse) Initialize() {
	for i, t := range p.tokens {
		if p.config.Unking && !p.grammar.Contains(t) {
			t = UnknownToken
		}

		for _, r := range p.grammar.Lexicon(t) {
			lexical := &Item{
				Span: Span{
					i: i,
					j: i + 1,
					n: r.Head,
				},
				weight: r.Weight(),
			}

			p.Push(lexical, 0)
		}
	}
}

func (p *parse) Combine(c1, c2 *Item, ri *grammar.NonLexical) *Item {
	return &Item{
		Span: Span{
			i: c1.i,
			j: c2.j,
			n: ri.Head,
		},
		weight:     c1.Weight() * c2.Weight() * ri.Weight(),
		backtracks: [2]*Item{c1, c2},
	}
}

func (p *parse) Chain(c *Item, ri *grammar.NonLexical) *Item {
	return &Item{
		Span: Span{
			i: c.i,
			j: c.j,
			n: ri.Head,
		},
		weight:     c.Weight() * ri.Weight(),
		backtracks: [2]*Item{c, nil},
	}
}

func (p *parse) Push(item *Item, threshold float64) bool {
	priority := p.ItemPriority(item)

	if threshold != 0 && priority < threshold {
		return false
	}

	p.queue.Push(item, priority)

	return true
}

func (p *parse) Tree(root *Item, tokens []string) (*tree.Tree, error) {
	var backtrack func(item *Item) (*tree.Tree, error)
	backtrack = func(item *Item) (*tree.Tree, error) {
		t := &tree.Tree{}

		if label, err := p.grammar.Symbols.Itoa(item.n); err != nil {
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
