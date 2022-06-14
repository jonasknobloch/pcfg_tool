package tool

import (
	"log"
	"os"
	"pcfg_tool/internal/grammar"
)

func Outside(rules, lexicon string, n string, output *os.File) {
	g := grammar.NewGrammar()

	g.SetInitial(n)

	if err := g.Import(rules, lexicon); err != nil {
		log.Fatal(err)
	}

	vs, err := grammar.NewViterbiScores(g)

	if err != nil {
		log.Fatal(err)
	}

	if err := vs.Export(output, g.Symbols); err != nil {
		log.Fatal(err)
	}
}
