package tool

import (
	"os"
	"pcfg_tool/internal/grammar"
	"pcfg_tool/internal/utility"
)

func Outside(rules, lexicon string, n string, outside string) error {
	g := grammar.NewGrammar()

	g.SetInitial(n)

	var r *os.File
	var l *os.File

	if f, err := os.Open(rules); err != nil {
		return err
	} else {
		r = f
		defer r.Close()
	}

	if f, err := os.Open(lexicon); err != nil {
		return err
	} else {
		l = f
		defer l.Close()
	}

	if err := g.Import(r, l); err != nil {
		return err
	}

	vs := grammar.NewViterbiScores()

	vs.CalcOutside(g)

	o, err := utility.CreateFile(outside)

	if err != nil {
		return err
	}

	defer o.Close()

	if err := vs.ExportOutside(o, g.Symbols); err != nil {
		return err
	}

	return nil
}
