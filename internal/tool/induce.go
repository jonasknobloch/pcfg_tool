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

		EvaluateTree(t, g)
	}

	g.Normalize()

	if err := g.SetInitial(n); err != nil {
		log.Fatal(err)
	}

	return g
}

func EvaluateTree(t *tree.Tree, g *grammar.Grammar) {
	t.Walk(func(t *tree.Tree) {
		if len(t.Children) == 0 {
			return
		}

		r, err := grammar.NewRule(t, g.Symbols)

		if err != nil {
			return
		}

		if err := g.AddRule(r); err != nil {
			panic(err) // TODO handle
		}
	})
}
