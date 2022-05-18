package pcfg

import "strconv"

func FormatWeight(weight float64) string {
	return strconv.FormatFloat(weight, 'f', -1, 64)
}

func NonLexicalToNonLexicalInt(r *NonLexical, st *SymbolTable) (*NonLexicalInt, error) {
	ri := &NonLexicalInt{}

	if head, err := st.Atoi(r.head); err != nil {
		return nil, err
	} else {
		ri.head = head
	}

	ri.body = make([]int, len(r.body))

	for i, b := range r.body {
		if j, err := st.Atoi(b); err != nil {
			return nil, err
		} else {
			ri.body[i] = j
		}
	}

	ri.rule = r

	return ri, nil
}
