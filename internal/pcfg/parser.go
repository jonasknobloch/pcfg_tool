package pcfg

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"golang.org/x/sync/semaphore"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
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
	grammar *Grammar
	rules   struct {
		value map[string][]Rule
		mutex sync.RWMutex
	}
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

	for r := range g.Weights() {
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
		grammar: g,
		rules: struct {
			value map[string][]Rule
			mutex sync.RWMutex
		}{value: rules, mutex: sync.RWMutex{}},
	}, nil
}

func (ps *Parser) Rules(body string) []Rule {
	ps.rules.mutex.Lock()
	defer ps.rules.mutex.Unlock()

	rules, ok := ps.rules.value[body]

	if !ok {
		rules = []Rule{}
	}

	return rules
}

func (ps *Parser) Parse(tokens []string) (*tree.Tree, bool) {
	p := &parse{
		tokens:  tokens,
		heap:    NewHeap(),
		matcher: NewMatcher(),
		parser:  ps,
	}

	return p.Parse()
}

func (ps *Parser) ParseFile(fs *bufio.Scanner) {
	ctx := context.TODO()
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))

	var wg sync.WaitGroup

	ph := NewPrintHeap()

	count := 0

	for fs.Scan() {
		text := fs.Text()

		if err := sem.Acquire(ctx, 1); err != nil {
			log.Fatalf("Failed to acquire semaphore: %v", err)
		}

		wg.Add(1)

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// 32 GB * 0.8 -> 256e8
		// 16 GB * 0.8 -> 128e8
		for m.Alloc > 256e8 {
			time.Sleep(100 * time.Millisecond)
		}

		go func(count int) {
			defer func() {
				defer sem.Release(1)
				defer wg.Done()
				defer ph.Print()
			}()

			tokens := strings.Split(text, " ")

			t, ok := ps.Parse(tokens)

			if !ok {
				ph.Push(&PrintJob{
					line:  fmt.Sprintf("(NOPARSE %s)", strings.Join(tokens, " ")),
					count: count,
				})
			} else {
				ph.Push(&PrintJob{
					line:  t.String(),
					count: count,
				})
			}
		}(count)

		count++
	}

	wg.Wait()
}

type parse struct {
	tokens  []string
	heap    *Heap
	matcher *Matcher
	parser  *Parser
}

func (p *parse) Parse() (*tree.Tree, bool) {
	p.Initialize()

	for !p.heap.Empty() {
		item, _ := p.heap.Pop()

		if ok := p.matcher.Add(item); !ok {
			continue
		}

		if item.n == p.parser.grammar.initial && item.i == 0 && item.j == len(p.tokens) {
			return p.Tree(item, p.tokens), true
		}

		for _, rule := range p.parser.Rules(item.n) {
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

func (p *parse) Initialize() {
	for i, t := range p.tokens {
		for _, r := range p.parser.Rules(t) {
			if _, ok := r.(*Lexical); !ok {
				continue
			}

			lexical := &Item{
				Span: Span{
					i: i,
					j: i + 1,
					n: r.Head(),
				},
				weight: p.parser.grammar.Weight(r),
			}

			p.heap.Push(lexical)
		}
	}
}

func (p *parse) Combine(c1, c2 *Item, r Rule) {
	i := &Item{
		Span: Span{
			i: c1.i,
			j: c2.j,
			n: r.Head(),
		},
		weight:     c1.Weight() * c2.Weight() * p.parser.grammar.Weight(r),
		backtracks: [2]*Item{c1, c2},
	}

	p.heap.Push(i)
}

func (p *parse) Chain(c *Item, r Rule) {
	i := &Item{
		Span: Span{
			i: c.i,
			j: c.j,
			n: r.Head(),
		},
		weight:     c.Weight() * p.parser.grammar.Weight(r),
		backtracks: [2]*Item{c, nil},
	}

	p.heap.Push(i)
}

func (p *parse) Tree(root *Item, tokens []string) *tree.Tree {
	var backtrack func(item *Item) *tree.Tree
	backtrack = func(item *Item) *tree.Tree {
		t := &tree.Tree{
			Label: item.n,
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

		if li == nil && ri == nil {
			t.Children = []*tree.Tree{
				{
					Label: tokens[item.i],
				},
			}
		}

		return t
	}

	return backtrack(root)
}
