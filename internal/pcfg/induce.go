package pcfg

import (
	"bufio"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"log"
	"os"
)

func Induce(file *os.File) *Grammar {
	dec := tree.NewDecoder()

	fs := bufio.NewScanner(file)

	fs.Split(bufio.ScanLines)

	g := NewGrammar()

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
	g.SetInitial(n)

	return g
}

func EvaluateTree(t *tree.Tree, g *Grammar) {
	t.Walk(func(t *tree.Tree) {
		if len(t.Children) == 0 {
			return
		}

		r, err := NewRule(t)

		if err != nil {
			return
		}

		g.AddRule(r, 1)
	})
}
