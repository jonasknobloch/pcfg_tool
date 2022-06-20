package tool

import (
	"bufio"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"pcfg_tool/internal/transform"
	"pcfg_tool/internal/utility"
)

func Transform(input, output string, callback func(t *tree.Tree)) error {
	var scanner *bufio.Scanner
	var writer *bufio.Writer

	if f, err := utility.OpenFile(input); err != nil {
		return err
	} else {
		scanner = bufio.NewScanner(f)

		defer f.Close()
	}

	if f, err := utility.CreateFile(output); err != nil {
		return err
	} else {
		writer = bufio.NewWriter(f)

		defer f.Close()
	}

	d := tree.NewDecoder()

	for scanner.Scan() {
		t := scanner.Text()

		tr, err := d.Decode(t)

		if err != nil {
			return err
		}

		callback(tr)

		if _, err := writer.WriteString(tr.String() + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func Markovize(horizontal, vertical int) func(*tree.Tree) {
	return func(t *tree.Tree) {
		transform.Markovize(t, horizontal, vertical)
	}
}

func Demarkovize() func(*tree.Tree) {
	return func(t *tree.Tree) {
		transform.Demarkovize(t)
	}
}
