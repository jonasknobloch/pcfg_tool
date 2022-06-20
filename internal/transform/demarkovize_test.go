package transform

import (
	"github.com/jonasknobloch/jinn/pkg/tree"
	"testing"
)

func TestDemarkovize(t *testing.T) {
	bt, _ := tree.NewDecoder().Decode("(ROOT (FRAG^<ROOT> (RB Not) (FRAG|<NP-TMP,.>^<ROOT> (NP-TMP^<FRAG,ROOT> (DT this) (NN year)) (. .))))")
	tr, _ := tree.NewDecoder().Decode("(ROOT (FRAG (RB Not) (NP-TMP (DT this) (NN year)) (. .)))")

	Demarkovize(bt)

	if !tr.Equals(bt) {
		t.Fatalf("expetced \n%s\n but got \n%s\n", tr.String(), bt.String())
	}
}
