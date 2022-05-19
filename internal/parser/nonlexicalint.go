package parser

import "pcfg_tool/internal/grammar"

type NonLexicalInt struct {
	Head int
	Body []int
	Rule *grammar.NonLexical
}

func NewNonLexicalInt(r *grammar.NonLexical, st *SymbolTable) (*NonLexicalInt, error) {
	ri := &NonLexicalInt{}

	if head, err := st.Atoi(r.Head); err != nil {
		return nil, err
	} else {
		ri.Head = head
	}

	ri.Body = make([]int, len(r.Body))

	for i, b := range r.Body {
		if j, err := st.Atoi(b); err != nil {
			return nil, err
		} else {
			ri.Body[i] = j
		}
	}

	ri.Rule = r

	return ri, nil
}
