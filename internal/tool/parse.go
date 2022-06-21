package tool

import (
	"bufio"
	"os"
	"pcfg_tool/internal/grammar"
	"pcfg_tool/internal/parser"
	"pcfg_tool/internal/utility"
)

func Parse(rules, lexicon string, n string, unking bool, path string, input string) error {
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

	var vs *grammar.ViterbiScores

	if path != "" {
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

	c := &parser.Config{
		Unking: unking,
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
