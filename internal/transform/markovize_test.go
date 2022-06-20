package transform

import (
	"github.com/jonasknobloch/jinn/pkg/tree"
	"testing"
)

func TestMarkovize(t *testing.T) {
	tr, _ := tree.NewDecoder().Decode("(ROOT (FRAG (RB Not) (NP-TMP (DT this) (NN year)) (. .)))")
	bt, _ := tree.NewDecoder().Decode("(ROOT (FRAG^<ROOT> (RB Not) (FRAG|<NP-TMP,.>^<ROOT> (NP-TMP^<FRAG,ROOT> (DT this) (NN year)) (. .))))")

	Markovize(tr, 999, 3)

	if !tr.Equals(bt) {
		t.Fatalf("expetced \n%s\n but got \n%s\n", tr.String(), bt.String())
	}
}
