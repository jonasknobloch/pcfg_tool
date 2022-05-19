package parser

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"golang.org/x/sync/semaphore"
	"log"
	"pcfg_tool/internal/config"
	"pcfg_tool/internal/grammar"
	"runtime"
	"strings"
	"sync"
	"time"
)

var ErrNoParse = errors.New("no parse")

type Parser struct {
	grammar *grammar.Grammar
	symbols *SymbolTable
	initial int
	lexicon struct {
		value map[string][]*grammar.Lexical
		mutex sync.RWMutex
	}
	rules struct {
		value map[int][]*NonLexicalInt
		mutex sync.RWMutex
	}
}

func NewParser(g *grammar.Grammar) (*Parser, error) {
	if g.Initial() == "" {
		return nil, errors.New("grammar initial not set")
	}

	p := &Parser{
		grammar: g,
		symbols: NewSymbolTable(),
	}

	p.lexicon = struct {
		value map[string][]*grammar.Lexical
		mutex sync.RWMutex
	}{value: make(map[string][]*grammar.Lexical), mutex: sync.RWMutex{}}

	p.rules = struct {
		value map[int][]*NonLexicalInt
		mutex sync.RWMutex
	}{value: make(map[int][]*NonLexicalInt), mutex: sync.RWMutex{}}

	add := func(k int, ir *NonLexicalInt) {
		if _, ok := p.rules.value[k]; !ok {
			p.rules.value[k] = make([]*NonLexicalInt, 0)
		}

		p.rules.value[k] = append(p.rules.value[k], ir)
	}

	if initial, err := p.symbols.Atoi(g.Initial()); err != nil {
		return nil, err
	} else {
		p.initial = initial
	}

	for r := range g.Weights() {
		switch v := r.(type) {
		case *grammar.Lexical:
			if _, ok := p.lexicon.value[v.KeyBody()]; !ok {
				p.lexicon.value[v.KeyBody()] = make([]*grammar.Lexical, 0)
			}

			p.lexicon.value[v.KeyBody()] = append(p.lexicon.value[v.KeyBody()], v)
		case *grammar.NonLexical:
			ri, err := NewNonLexicalInt(v, p.symbols)

			if err != nil {
				return nil, err
			}

			add(ri.Body[0], ri)

			if ri.Body[len(ri.Body)-1] != ri.Body[0] {
				add(ri.Body[len(ri.Body)-1], ri)
			}
		default:
			panic("unknown rule type")
		}
	}

	return p, nil
}

func (ps *Parser) Lexicon(body string) []*grammar.Lexical {
	ps.lexicon.mutex.Lock()
	defer ps.lexicon.mutex.Unlock()

	lexicon, ok := ps.lexicon.value[body]

	if !ok {
		lexicon = []*grammar.Lexical{}
	}

	return lexicon
}

func (ps *Parser) Rules(body int) []*NonLexicalInt {
	ps.rules.mutex.Lock()
	defer ps.rules.mutex.Unlock()

	rules, ok := ps.rules.value[body]

	if !ok {
		rules = []*NonLexicalInt{}
	}

	return rules
}

func (ps *Parser) Parse(tokens []string) (*tree.Tree, error) {
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
	sem := semaphore.NewWeighted(config.Config.WorkerPoolSize)

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

		for m.Alloc > config.Config.AllocThreshold {
			time.Sleep(config.Config.ReadMemStatsRate)
			runtime.ReadMemStats(&m)
		}

		go func(count int) {
			defer func() {
				defer sem.Release(1)
				defer wg.Done()
				defer ph.Print()
			}()

			tokens := strings.Split(text, " ")

			t, err := ps.Parse(tokens)

			if err == nil {
				ph.Push(&PrintJob{
					line:  t.String(),
					count: count,
				})
			} else if err == ErrNoParse {
				ph.Push(&PrintJob{
					line:  fmt.Sprintf("(NOPARSE %s)", text),
					count: count,
				})
			} else {
				log.Fatal(err)
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

func (p *parse) Parse() (*tree.Tree, error) {
	if err := p.Initialize(); err != nil {
		return nil, err
	}

	for !p.heap.Empty() {
		item, _ := p.heap.Pop()

		if ok := p.matcher.Add(item); !ok {
			continue
		}

		if item.n == p.parser.initial && item.i == 0 && item.j == len(p.tokens) {
			return p.Tree(item, p.tokens)
		}

		for _, rule := range p.parser.Rules(item.n) {
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
		for _, r := range p.parser.Lexicon(t) {
			n, err := p.parser.symbols.Atoi(r.KeyHead())

			if err != nil {
				return err
			}

			lexical := &Item{
				Span: Span{
					i: i,
					j: i + 1,
					n: n,
				},
				weight: p.parser.grammar.Weight(r),
			}

			p.heap.Push(lexical)
		}
	}

	return nil
}

func (p *parse) Combine(c1, c2 *Item, ri *NonLexicalInt) {
	i := &Item{
		Span: Span{
			i: c1.i,
			j: c2.j,
			n: ri.Head,
		},
		weight:     c1.Weight() * c2.Weight() * p.parser.grammar.Weight(ri.Rule),
		backtracks: [2]*Item{c1, c2},
	}

	p.heap.Push(i)
}

func (p *parse) Chain(c *Item, ri *NonLexicalInt) {
	i := &Item{
		Span: Span{
			i: c.i,
			j: c.j,
			n: ri.Head,
		},
		weight:     c.Weight() * p.parser.grammar.Weight(ri.Rule),
		backtracks: [2]*Item{c, nil},
	}

	p.heap.Push(i)
}

func (p *parse) Tree(root *Item, tokens []string) (*tree.Tree, error) {
	var backtrack func(item *Item) (*tree.Tree, error)
	backtrack = func(item *Item) (*tree.Tree, error) {
		t := &tree.Tree{}

		if label, ok := p.parser.symbols.Itoa(item.n); !ok {
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
