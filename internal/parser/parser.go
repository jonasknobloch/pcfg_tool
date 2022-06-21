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
	viterbi *grammar.ViterbiScores
	config  *Config
}

func NewParser(g *grammar.Grammar, vs *grammar.ViterbiScores, c *Config) (*Parser, error) {
	if g.Initial() == 0 {
		return nil, errors.New("grammar initial not set")
	}

	return &Parser{
		grammar: g,
		viterbi: vs,
		config:  c,
	}, nil
}

func (ps *Parser) Parse(tokens []string) (*tree.Tree, error) {
	p := &parse{
		tokens:  tokens,
		heap:    NewHeap(),
		matcher: NewMatcher(),
		grammar: ps.grammar,
		viterbi: ps.viterbi,
		config:  ps.config,
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
