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
		text := fs.Text()
		tokens := strings.Split(text, " ")

		t, err := p.Parse(tokens)

		if err == nil {
			fmt.Println(t)
		} else if err == ErrNoParse {
			fmt.Printf("(NOPARSE %s)\n", text)
		} else {
			log.Fatal(err)
		}
	}
}
