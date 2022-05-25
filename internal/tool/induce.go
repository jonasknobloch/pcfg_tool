package tool

import (
	"bufio"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"log"
	"os"
	"pcfg_tool/internal/grammar"
)

func Induce(file *os.File) *grammar.Grammar {
	dec := tree.NewDecoder()

	fs := bufio.NewScanner(file)

	fs.Split(bufio.ScanLines)

	g := grammar.NewGrammar()

	var n string

	for fs.Scan() {
		gold := fs.Text()

		t, err := dec.Decode(gold)

		if err != nil {
			log.Fatal(err)
		}

		if n != t.Label {
			n = t.Label
		}

		if err := EvaluateTree(t, g); err != nil {
			log.Fatal(err)
		}
	}

	g.Normalize()

	if err := g.SetInitial(n); err != nil {
		log.Fatal(err)
	}

	return g
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

		r, k, err := grammar.NewRule(t, g.Symbols)

		if err != nil {
			return err
		}

		if err := g.AddRule(r, k); err != nil {
			return err
		}

		return nil
	})
}
