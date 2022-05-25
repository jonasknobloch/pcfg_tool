package grammar

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"pcfg_tool/internal/utility"
	"strconv"
	"strings"
	"sync"
)

type Grammar struct {
	initial int
	rules   struct {
		left  map[int][]*NonLexical
		right map[int][]*NonLexical
		key   map[uint64]*NonLexical
		mutex sync.RWMutex
	}
	lexicon struct {
		left  map[int][]*Lexical
		right map[string][]*Lexical
		key   map[uint64]*Lexical
		mutex sync.RWMutex
	}
	words   map[string]struct{}
	Symbols *SymbolTable
}

var ErrUnknownRuleType = errors.New("unknown rule type")

func NewGrammar() *Grammar {
	return &Grammar{
		rules: struct {
			left  map[int][]*NonLexical
			right map[int][]*NonLexical
			key   map[uint64]*NonLexical
			mutex sync.RWMutex
		}{
			left:  make(map[int][]*NonLexical),
			right: make(map[int][]*NonLexical),
			key:   make(map[uint64]*NonLexical),
			mutex: sync.RWMutex{},
		},
		lexicon: struct {
			left  map[int][]*Lexical
			right map[string][]*Lexical
			key   map[uint64]*Lexical
			mutex sync.RWMutex
		}{
			left:  make(map[int][]*Lexical),
			right: make(map[string][]*Lexical),
			key:   make(map[uint64]*Lexical),
			mutex: sync.RWMutex{},
		},
		words:   make(map[string]struct{}),
		Symbols: NewSymbolTable(),
	}
}

func (g *Grammar) Initial() int {
	return g.initial
}

func (g *Grammar) SetInitial(n string) error {
	var err error

	g.initial, err = g.Symbols.Atoi(n)

	return err
}

func (g *Grammar) AddRule(rule Rule) error {
	var err error

	switch v := rule.(type) {
	case *NonLexical:
		err = g.AddNonLexical(v)
	case *Lexical:
		err = g.AddLexical(v)
	default:
		err = ErrUnknownRuleType
	}

	return err
}

func (g *Grammar) AddNonLexical(nonLexical *NonLexical) error {
	if nl, ok := g.rules.key[nonLexical.key]; !ok {
		g.rules.key[nonLexical.key] = nonLexical
	} else {
		nl.weight += nonLexical.weight
		return nil
	}

	if _, ok := g.rules.left[nonLexical.Head]; !ok {
		g.rules.left[nonLexical.Head] = make([]*NonLexical, 0)
	}

	g.rules.left[nonLexical.Head] = append(g.rules.left[nonLexical.Head], nonLexical)

	for _, b := range nonLexical.Body {
		if _, ok := g.rules.right[b]; !ok {
			g.rules.right[b] = make([]*NonLexical, 0)
		}

		g.rules.right[b] = append(g.rules.right[b], nonLexical)
	}

	return nil
}

func (g *Grammar) AddLexical(lexical *Lexical) error {
	if l, ok := g.lexicon.key[lexical.key]; !ok {
		g.lexicon.key[lexical.key] = lexical
	} else {
		l.weight += lexical.weight
		return nil
	}

	if _, ok := g.lexicon.left[lexical.Head]; !ok {
		g.lexicon.left[lexical.Head] = make([]*Lexical, 0)
	}

	g.lexicon.left[lexical.Head] = append(g.lexicon.left[lexical.Head], lexical)

	if _, ok := g.lexicon.right[lexical.Body]; !ok {
		g.lexicon.right[lexical.Body] = make([]*Lexical, 0)
	}

	g.lexicon.right[lexical.Body] = append(g.lexicon.right[lexical.Body], lexical)

	g.words[lexical.Body] = struct{}{}

	return nil
}

func (g *Grammar) Normalize() {
	symbols := make(map[int]float64)

	for s, nls := range g.rules.left {
		for _, nl := range nls {
			symbols[s] += nl.weight
		}
	}

	for s, ls := range g.lexicon.left {
		for _, l := range ls {
			symbols[s] += l.weight
		}
	}

	for s, nls := range g.rules.left {
		for _, nl := range nls {
			nl.weight /= symbols[s]
		}
	}

	for s, ls := range g.lexicon.left {
		for _, l := range ls {
			l.weight /= symbols[s]
		}
	}
}

func (g *Grammar) IsNormalized() bool {
	symbols := make(map[int]float64)

	for s, nls := range g.rules.left {
		for _, nl := range nls {
			symbols[s] += nl.weight
		}
	}

	for s, ls := range g.lexicon.left {
		for _, l := range ls {
			symbols[s] += l.weight
		}
	}

	for _, w := range symbols {
		if w-0.1 > 1 && w+0.1 < 1 {
			return false
		}
	}

	return true
}

func (g *Grammar) Rules(body int) []*NonLexical {
	g.rules.mutex.Lock()
	defer g.rules.mutex.Unlock()

	rules, ok := g.rules.right[body]

	if !ok {
		rules = []*NonLexical{}
	}

	return rules
}

func (g *Grammar) Lexicon(body string) []*Lexical {
	g.lexicon.mutex.Lock()
	defer g.lexicon.mutex.Unlock()

	lexicon, ok := g.lexicon.right[body]

	if !ok {
		lexicon = []*Lexical{}
	}

	return lexicon
}

func (g *Grammar) Print() error {
	printRule := func(rule Rule) error {
		var sb strings.Builder

		key, err := rule.String(g.Symbols)

		if err != nil {
			return err
		}

		sb.WriteString(key)
		sb.WriteString(fmt.Sprintf(" %s", utility.FormatWeight(rule.Weight())))

		fmt.Println(sb.String())

		return nil
	}

	for _, nls := range g.rules.left {
		for _, nl := range nls {
			err := printRule(nl)

			if err != nil {
				return err
			}
		}
	}

	for _, ls := range g.lexicon.left {
		for _, l := range ls {
			err := printRule(l)

			if err != nil {
				return err
			}
		}
	}

	for t := range g.words {
		fmt.Println(t)
	}

	return nil
}

func (g *Grammar) Import(rules, lexicon string) error {
	var rS *bufio.Scanner
	var lS *bufio.Scanner

	if file, err := os.Open(rules); err != nil {
		return err
	} else {
		rS = bufio.NewScanner(file)
		defer file.Close()
	}

	if file, err := os.Open(lexicon); err != nil {
		return err
	} else {
		lS = bufio.NewScanner(file)
		defer file.Close()
	}

	for rS.Scan() {
		t := strings.Split(rS.Text(), " ")

		r, err := NewNonLexical(t[0], t[2:len(t)-1], g.Symbols)

		if err != nil {
			return err
		}

		w, err := strconv.ParseFloat(t[len(t)-1], 64)

		if err != nil {
			return err
		}

		r.weight = w

		if err := g.AddRule(r); err != nil {
			return err
		}
	}

	for lS.Scan() {
		t := strings.Split(lS.Text(), " ")

		r, err := NewLexical(t[0], t[1], g.Symbols)

		if err != nil {
			return err
		}

		w, err := strconv.ParseFloat(t[2], 64)

		if err != nil {
			return err
		}

		r.weight = w

		if err := g.AddRule(r); err != nil {
			return err
		}
	}

	return nil
}

func (g *Grammar) Export(grammar string) error {
	var rules *bufio.Writer
	var lexicon *bufio.Writer
	var words *bufio.Writer

	if file, err := os.Create(grammar + ".rules"); err != nil {
		return err
	} else {
		rules = bufio.NewWriter(file)
	}

	if file, err := os.Create(grammar + ".lexicon"); err != nil {
		return err
	} else {
		lexicon = bufio.NewWriter(file)
	}

	if file, err := os.Create(grammar + ".words"); err != nil {
		return err
	} else {
		words = bufio.NewWriter(file)
	}

	for _, nls := range g.rules.left {
		for _, nl := range nls {
			s, err := nl.String(g.Symbols)

			if err != nil {
				return err
			}

			if _, err := rules.WriteString(s); err != nil {
				return err
			}
		}
	}

	for _, ls := range g.lexicon.left {
		for _, l := range ls {
			s, err := l.String(g.Symbols)

			if err != nil {
				return err
			}

			if _, err := lexicon.WriteString(s); err != nil {
				return err
			}
		}
	}

	for t := range g.words {
		_, err := words.WriteString(t + "\n")

		if err != nil {
			return err
		}
	}

	if err := rules.Flush(); err != nil {
		return err
	}

	if err := lexicon.Flush(); err != nil {
		return err
	}

	if err := words.Flush(); err != nil {
		return err
	}

	return nil
}
