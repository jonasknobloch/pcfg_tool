package pcfg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Grammar struct {
	initial   string
	weights   map[*Rule]float64
	rules     map[string]map[string]*Rule
	terminals map[string]struct{}
}

func NewGrammar() *Grammar {
	return &Grammar{
		initial:   "ROOT", // TODO from flag?
		weights:   make(map[*Rule]float64),
		rules:     make(map[string]map[string]*Rule),
		terminals: make(map[string]struct{}),
	}
}

func (g *Grammar) AddRule(rule *Rule, weight float64) {
	head := (*rule).Head()
	body := (*rule).Body()

	if _, ok := g.rules[head]; !ok {
		g.rules[head] = make(map[string]*Rule, 0)
	}

	if r, ok := g.rules[head][body]; ok {
		rule = r
	} else {
		g.rules[head][body] = rule
	}

	if _, ok := (*rule).(*Lexical); ok {
		g.terminals[body] = struct{}{}
	}

	g.weights[rule] += weight
}

func (g *Grammar) Normalize() {
	for _, bodies := range g.rules {
		sum := float64(0)

		for _, rule := range bodies {
			sum += g.weights[rule]
		}

		for _, rule := range bodies {
			g.weights[rule] = g.weights[rule] / sum
		}
	}
}

func (g *Grammar) IsNormalized() bool {
	for _, bodies := range g.rules {
		sum := float64(0)

		for _, rule := range bodies {
			sum += g.weights[rule]
		}

		if sum-0.1 > 1 && sum+0.1 < 1 {
			return false
		}
	}

	return true
}

func (g *Grammar) Print() {
	for r, w := range g.weights {
		var sb strings.Builder

		sb.WriteString((*r).String())
		sb.WriteString(fmt.Sprintf(" %s", FormatWeight(w)))

		fmt.Println(sb.String())
	}

	for t := range g.terminals {
		fmt.Println(t)
	}
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

	for r, w := range g.weights {
		var sb strings.Builder

		sb.WriteString((*r).String())
		sb.WriteString(fmt.Sprintf(" %s\n", FormatWeight(w)))

		var err error

		switch (*r).(type) {
		case *NonLexical:
			_, err = rules.WriteString(sb.String())
		case *Lexical:
			_, err = lexicon.WriteString(sb.String())
		default:
			panic("unknown rule type")
		}

		if err != nil {
			return err
		}
	}

	for t := range g.terminals {
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
