package grammar

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Grammar struct {
	initial int
	rules   struct {
		left  map[int][]*NonLexical
		right map[int][]*NonLexical
		body  map[[2]int][]*NonLexical
		key   map[string]*NonLexical
	}
	lexicon struct {
		left  map[int][]*Lexical
		right map[string][]*Lexical
		key   map[string]*Lexical
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
			body  map[[2]int][]*NonLexical
			key   map[string]*NonLexical
		}{
			left:  make(map[int][]*NonLexical),
			right: make(map[int][]*NonLexical),
			body:  make(map[[2]int][]*NonLexical),
			key:   make(map[string]*NonLexical),
		},
		lexicon: struct {
			left  map[int][]*Lexical
			right map[string][]*Lexical
			key   map[string]*Lexical
		}{
			left:  make(map[int][]*Lexical),
			right: make(map[string][]*Lexical),
			key:   make(map[string]*Lexical),
		},
		words:   make(map[string]struct{}),
		Symbols: NewSymbolTable(),
	}
}

func (g *Grammar) Initial() int {
	return g.initial
}

func (g *Grammar) SetInitial(n string) {
	g.initial = g.Symbols.Atoi(n)
}

func (g *Grammar) AddRule(rule Rule) error {
	switch v := rule.(type) {
	case *NonLexical:
		g.AddNonLexical(v)
	case *Lexical:
		g.AddLexical(v)
	default:
		return ErrUnknownRuleType
	}

	return nil
}

func (g *Grammar) AddNonLexical(nonLexical *NonLexical) {
	key := nonLexical.Key()

	if nl, ok := g.rules.key[key]; !ok {
		g.rules.key[key] = nonLexical
	} else {
		nl.weight += nonLexical.weight
		return
	}

	body := [2]int{nonLexical.Body[0]}

	if l := len(nonLexical.Body); l > 1 {
		body[1] = nonLexical.Body[l-1]
	}

	if _, ok := g.rules.body[body]; !ok {
		g.rules.body[body] = make([]*NonLexical, 0)
	}

	g.rules.body[body] = append(g.rules.body[body], nonLexical)

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
}

func (g *Grammar) AddLexical(lexical *Lexical) {
	key := lexical.Key()

	if l, ok := g.lexicon.key[key]; !ok {
		g.lexicon.key[key] = lexical
	} else {
		l.weight += lexical.weight
		return
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
		if w > 1.1 || w < 0.9 {
			return false
		}
	}

	return true
}

func (g *Grammar) Rules(body int) []*NonLexical {
	rules, ok := g.rules.right[body]

	if !ok {
		rules = []*NonLexical{}
	}

	return rules
}

func (g *Grammar) ExactRules(body [2]int) []*NonLexical {
	rules, ok := g.rules.body[body]

	if !ok {
		rules = []*NonLexical{}
	}

	return rules
}

func (g *Grammar) Lexicon(body string) []*Lexical {
	lexicon, ok := g.lexicon.right[body]

	if !ok {
		lexicon = []*Lexical{}
	}

	return lexicon
}

func (g *Grammar) Import(rules, lexicon *os.File) error {
	rS := bufio.NewScanner(rules)
	lS := bufio.NewScanner(lexicon)

	for rS.Scan() {
		t := strings.Split(rS.Text(), " ")

		w, err := strconv.ParseFloat(t[len(t)-1], 64)

		if err != nil {
			return err
		}

		g.AddNonLexical(NewNonLexical(t[0], t[2:len(t)-1], w, g.Symbols))
	}

	for lS.Scan() {
		t := strings.Split(lS.Text(), " ")

		w, err := strconv.ParseFloat(t[2], 64)

		if err != nil {
			return err
		}

		g.AddLexical(NewLexical(t[0], t[1], w, g.Symbols))
	}

	return nil
}

func (g *Grammar) Export(rules, lexicon, words *os.File) error {
	rW := bufio.NewWriter(rules)
	lW := bufio.NewWriter(lexicon)
	wW := bufio.NewWriter(words)

	for _, nls := range g.rules.left {
		for _, nl := range nls {
			s, err := nl.String(g.Symbols)

			if err != nil {
				return err
			}

			if _, err := rW.WriteString(s + "\n"); err != nil {
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

			if _, err := lW.WriteString(s + "\n"); err != nil {
				return err
			}
		}
	}

	for t := range g.words {
		_, err := wW.WriteString(t + "\n")

		if err != nil {
			return err
		}
	}

	if err := rW.Flush(); err != nil {
		return err
	}

	if err := lW.Flush(); err != nil {
		return err
	}

	if err := wW.Flush(); err != nil {
		return err
	}

	return nil
}
