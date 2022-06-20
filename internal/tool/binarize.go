package tool

import (
	"bufio"
	"github.com/jonasknobloch/jinn/pkg/tree"
	"pcfg_tool/internal/transform"
	"pcfg_tool/internal/utility"
)

func Binarize(input, output string, horizontal, vertical int) error {
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

		transform.Markovize(tr, horizontal, vertical)

		if _, err := writer.WriteString(tr.String() + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}
