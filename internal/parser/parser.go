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
