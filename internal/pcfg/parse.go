package pcfg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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

	fs.Split(bufio.ScanLines)

	for fs.Scan() {
		tokens := strings.Split(fs.Text(), " ")

		t, ok := p.Parse(tokens)

		if !ok {
			fmt.Printf("(NOPARSE %s)\n", strings.Join(p.tokens, " "))
		} else {
			fmt.Println(t)
		}
	}
}
