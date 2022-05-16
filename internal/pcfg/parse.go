package pcfg

import (
	"bufio"
	"log"
	"os"
)

func Parse(rules, lexicon string, n string, file *os.File) {
	g := NewGrammar()

	g.SetInitial(n)

	if err := g.Import(rules, lexicon); err != nil {
		log.Fatal(err)
	}

	p, err := NewParser(g)

	if err != nil {
		log.Fatal(err)
	}

	fs := bufio.NewScanner(file)

	p.ParseFile(fs)
}
