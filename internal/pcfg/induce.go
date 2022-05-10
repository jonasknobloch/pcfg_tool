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

	for fs.Scan() {
		gold := fs.Text()

		t, err := dec.Decode(gold)

		if err != nil {
			log.Fatal(err)
		}

		EvaluateTree(t, g)
	}

	g.Normalize()

	return g
}

func EvaluateTree(t *tree.Tree, g *Grammar) {
	t.Walk(func(t *tree.Tree) {
		if len(t.Children) == 0 {
			return
		}

		r := NewRule(t)

		g.AddRule(r, 1)
	})
}
