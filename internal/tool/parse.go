package tool

import (
	"bufio"
	"log"
	"os"
	"pcfg_tool/internal/grammar"
	"pcfg_tool/internal/parser"
)

func Parse(rules, lexicon string, n string, file *os.File) {
	g := grammar.NewGrammar()

	g.SetInitial(n)

	if err := g.Import(rules, lexicon); err != nil {
		log.Fatal(err)
	}

	p, err := parser.NewParser(g)

	if err != nil {
		log.Fatal(err)
	}

	fs := bufio.NewScanner(file)

	p.ParseFile(fs)
}
