package grammar

import (
	"pcfg_tool/internal/utility"
	"strings"
)

type NonLexical struct {
	Head   int
	Body   []int
	weight float64
}

func NewNonLexical(head string, body []string, symbols *SymbolTable) (*NonLexical, string) {
	if len(body) == 0 {
		body = []string{"X"}
	}

	nl := &NonLexical{
		weight: 1,
		Body:   make([]int, len(body)),
	}

	nl.Head = symbols.Atoi(head)

	for i, b := range body {
		nl.Body[i] = symbols.Atoi(b)
	}

	key := head + " " + strings.Join(body, " ")

	return nl, key
}

func (nl *NonLexical) Weight() float64 {
	return nl.weight
}

func (nl *NonLexical) String(st *SymbolTable) (string, error) {
	var sb strings.Builder

	if head, err := st.Itoa(nl.Head); err != nil {
		return "", err
	} else {
		sb.WriteString(head)
	}

	sb.WriteString(" ->")

	for _, b := range nl.Body {
		v, err := st.Itoa(b)

		if err != nil {
			return "", err
		}

		sb.WriteString(" ")
		sb.WriteString(v)
	}

	sb.WriteString(" ")
	sb.WriteString(utility.FormatWeight(nl.weight))

	return sb.String(), nil
}
