package tool

import (
	"bufio"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"os"
	"pcfg_tool/internal/grammar"
	"pcfg_tool/internal/utility"
)

func Induce(input string, rules, lexicon, words string) error {
	var scanner *bufio.Scanner

	if f, err := utility.OpenFile(input); err != nil {
		return err
	} else {
		scanner = bufio.NewScanner(f)

		defer f.Close()
	}

	d := tree.NewDecoder()
	g := grammar.NewGrammar()

	var n string

	for scanner.Scan() {
		gold := scanner.Text()

		t, err := d.Decode(gold)

		if err != nil {
			return err
		}

		if n != t.Label {
			n = t.Label
		}

		if err := EvaluateTree(t, g); err != nil {
			return err
		}
	}

	g.Normalize()
	g.SetInitial(n)

	var r *os.File
	var l *os.File
	var w *os.File

	var err error

	if r, err = utility.CreateFile(rules); err != nil {
		return err
	}

	defer r.Close()

	if l, err = utility.CreateFile(lexicon); err != nil {
		return err
	}

	defer l.Close()

	if w, err = utility.CreateFile(words); err != nil {
		return err
	}

	defer w.Close()

	return g.Export(r, l, w)
}

func EvaluateTree(t *tree.Tree, g *grammar.Grammar) error {
	var walk func(*tree.Tree, func(t *tree.Tree) error) error
	walk = func(t *tree.Tree, cb func(t *tree.Tree) error) error {
		if err := cb(t); err != nil {
			return err
		}

		for _, c := range t.Children {
			if err := walk(c, cb); err != nil {
				return err
			}
		}

		return nil
	}

	return walk(t, func(t *tree.Tree) error {
		if len(t.Children) == 0 {
			return nil
		}

		r, err := grammar.NewRule(t, g.Symbols)

		if err != nil {
			return err
		}

		if err := g.AddRule(r); err != nil {
			return err
		}

		return nil
	})
}
