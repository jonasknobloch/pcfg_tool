package grammar

import (
	"bufio"
	"os"
	"pcfg_tool/internal/utility"
	"strconv"
	"strings"
)

type ViterbiScores struct {
	inside  map[int]float64
	outside map[int]float64
}

func NewViterbiScores() *ViterbiScores {
	vs := &ViterbiScores{
		inside:  make(map[int]float64),
		outside: make(map[int]float64),
	}

	return vs
}

func (vs *ViterbiScores) Outside(n int) float64 {
	return vs.outside[n]
}

func (vs *ViterbiScores) calcInside(g *Grammar) {
	for head, rules := range g.lexicon.left {
		for _, lexical := range rules {
			if vs.inside[head] < lexical.weight {
				vs.inside[head] = lexical.weight
			}
		}
	}

	converged := false

	for !converged {
		converged = true

		for head, rules := range g.rules.left {
			for _, nonLexical := range rules {
				weight := nonLexical.weight

				weight *= vs.inside[nonLexical.Body[0]]

				if len(nonLexical.Body) == 2 {
					weight *= vs.inside[nonLexical.Body[1]]
				}

				if weight > vs.inside[head] {
					vs.inside[head] = weight
					converged = false
				}
			}
		}
	}
}

func (vs *ViterbiScores) CalcOutside(g *Grammar) {
	vs.calcInside(g)

	vs.outside[g.initial] = 1

	converged := false

	for !converged {
		converged = true

		for body, rules := range g.rules.right {
			for _, nonLexical := range rules {
				weight := nonLexical.weight

				weight *= vs.outside[nonLexical.Head]

				if len(nonLexical.Body) == 2 {
					if nonLexical.Body[0] != body {
						weight *= vs.inside[nonLexical.Body[0]]
					} else {
						weight *= vs.inside[nonLexical.Body[1]]
					}
				}

				if weight > vs.outside[body] {
					vs.outside[body] = weight
					converged = false
				}
			}
		}
	}
}

func (vs *ViterbiScores) ImportOutside(outside *os.File, symbols *SymbolTable) error {
	scanner := bufio.NewScanner(outside)

	for scanner.Scan() {
		t := strings.Split(scanner.Text(), " ")

		w, err := strconv.ParseFloat(t[1], 64)

		if err != nil {
			return err
		}

		vs.outside[symbols.Atoi(t[0])] = w
	}

	return nil
}

func (vs *ViterbiScores) ExportOutside(outside *os.File, symbols *SymbolTable) error {
	writer := bufio.NewWriter(outside)

	for v, w := range vs.outside {
		s, err := symbols.Itoa(v)

		if err != nil {
			return err
		}

		if _, err := writer.WriteString(s + " " + utility.FormatWeight(w) + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}
