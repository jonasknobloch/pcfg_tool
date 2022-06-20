package transform

import (
	"github.com/jonasknobloch/jinn/pkg/tree"
	"strings"
)

func Markovize(t *tree.Tree, h, v int) {
	var walk func(*tree.Tree, []*tree.Tree, func(*tree.Tree, []*tree.Tree))
	walk = func(t *tree.Tree, ps []*tree.Tree, cb func(*tree.Tree, []*tree.Tree)) {
		cb(t, ps)

		nps := make([]*tree.Tree, len(ps), len(ps)+1)
		copy(nps, ps)
		nps = append(nps, t)

		for _, c := range t.Children {
			walk(c, nps, cb)
		}
	}

	joinLabels := func(prefix string, reverse bool, ts []*tree.Tree) string {
		var sb strings.Builder

		if len(ts) > 0 {
			sb.WriteString(prefix)
			sb.WriteString("<")

			for i := 0; i < len(ts); i++ {
				var j int

				if !reverse {
					j = i
				} else {
					j = len(ts) - 1 - i
				}

				if i > 0 {
					sb.WriteString(",")
				}

				sb.WriteString(ts[j].Label)
			}

			sb.WriteString(">")
		}

		return sb.String()
	}

	rs := make(map[*tree.Tree]string) // label replacements
	as := make(map[*tree.Tree]string) // label additions

	walk(t, []*tree.Tree{}, func(t *tree.Tree, ps []*tree.Tree) {
		if len(t.Children) == 0 {
			return
		}

		if len(t.Children) == 1 && len(t.Children[0].Children) == 0 {
			return
		}

		i := len(ps) - v
		j := len(ps)

		if i < 0 {
			i = 0
		}

		as[t] = joinLabels("^", true, ps[i:j])
	})

	walk(t, []*tree.Tree{}, func(t *tree.Tree, ps []*tree.Tree) {
		if len(t.Children) == 0 {
			return
		}

		if len(t.Children) == 1 && len(t.Children[0].Children) == 0 {
			return
		}

		if len(t.Children) > 2 {
			head := make([]*tree.Tree, 1, 2)
			tail := make([]*tree.Tree, len(t.Children)-1)

			copy(head, t.Children[0:1])
			copy(tail, t.Children[1:])

			i := 0
			j := len(tail)

			if j > h {
				j = h
			}

			r := t.Label + joinLabels("|", false, tail[i:j])

			st := &tree.Tree{
				Label:    t.Label,
				Children: tail,
			}

			head = append(head, st)

			t.Children = head

			rs[st] = r
			as[st] = as[t]
		}
	})

	for st, r := range rs {
		st.Label = r
	}

	for st, a := range as {
		st.Label = st.Label + a
	}
}
