package pcfg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Parse(rules, lexicon string, file *os.File) {
	g := NewGrammar()

	if err := g.Import(rules, lexicon); err != nil {
		log.Fatal(err)
	}

	p := NewParser(g)

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
