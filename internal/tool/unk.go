package tool

import (
	"github.com/jonasknobloch/jinn/pkg/tree"
	"os"
	"pcfg_tool/internal/parser"
	"pcfg_tool/internal/utility"
)

func Unk(input, output string, threshold int) error {
	w := make(map[string]int)

	count := func(t *tree.Tree) {
		for _, l := range t.Leaves() {
			w[l.Label]++
		}
	}

	replace := func(t *tree.Tree) {
		for _, l := range t.Leaves() {
			if w[l.Label] <= threshold {
				l.Label = parser.UnknownToken
			}
		}
	}

	var temp *os.File
	var name string

	{
		var file *os.File
		var err error

		file, err = utility.OpenFile(input)

		if err != nil {
			return err
		}

		defer file.Close()

		temp, err = utility.CopyToTemp(file)

		if err != nil {
			return err
		}

		name = temp.Name()

		defer temp.Close()
		defer os.Remove(name)
	}

	if err := Transform(name, "/dev/null", count); err != nil {
		return err
	}

	if err := Transform(name, output, replace); err != nil {
		return err
	}

	return nil
}
