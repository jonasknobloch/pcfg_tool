package tool

import (
	"bufio"
	"os"
	"pcfg_tool/internal/grammar"
	"pcfg_tool/internal/parser"
	"pcfg_tool/internal/utility"
)

func Parse(rules, lexicon string, n string, unking bool, threshold float64, rank int, path string, input string) error {
	g := grammar.NewGrammar()

	g.SetInitial(n)

	var r *os.File
	var l *os.File

	if f, err := os.Open(rules); err != nil {
		return err
	} else {
		r = f
	}

	if f, err := os.Open(lexicon); err != nil {
		return err
	} else {
		l = f
	}

	if err := g.Import(r, l); err != nil {
		return err
	}

	c := &parser.Config{
		Unking: unking,
	}

	if path != "" {
		c.AStar = true
	}

	if threshold != 0 || rank != 0 {
		c.Prune = true
		c.Threshold = threshold
		c.Rank = rank
	}

	var vs *grammar.ViterbiScores

	if c.AStar {
		var o *os.File

		if f, err := os.Open(path); err != nil {
			return err
		} else {
			o = f
		}

		vs = grammar.NewViterbiScores()

		if err := vs.ImportOutside(o, g.Symbols); err != nil {
			return err
		}
	}

	p, err := parser.NewParser(g, vs, c)

	if err != nil {
		return err
	}

	var scanner *bufio.Scanner

	if f, err := utility.OpenFile(input); err != nil {
		return err
	} else {
		scanner = bufio.NewScanner(f)
	}

	p.ParseFile(scanner)

	return nil
}
