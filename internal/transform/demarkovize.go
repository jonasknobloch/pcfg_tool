package transform

import (
	"github.com/jonasknobloch/jinn/pkg/tree"
	"strings"
)

func Demarkovize(t *tree.Tree) {
	t.Walk(func(st *tree.Tree) {
		if i := strings.Index(st.Label, "<"); i != -1 {
			st.Label = st.Label[:i-1]
		}

		converged := false

		for !converged {
			converged = true

			for i, c := range st.Children {
				if !strings.Contains(c.Label, "|") {
					continue
				}

				converged = false

				st.Children = st.Children[:i]

				for _, cc := range c.Children {
					st.Children = append(st.Children, cc)
				}
			}
		}
	})
}
