package grammar

import (
	"pcfg_tool/internal/utility"
	"strconv"
	"strings"
)

type NonLexical struct {
	Head   int
	Body   []int
	weight float64
}

func NewNonLexical(head string, body []string, weight float64, symbols *SymbolTable) *NonLexical {
	if len(body) == 0 {
		body = []string{"X"}
	}

	nl := &NonLexical{
		weight: weight,
		Body:   make([]int, len(body)),
	}

	nl.Head = symbols.Atoi(head)

	for i, b := range body {
		nl.Body[i] = symbols.Atoi(b)
	}

	return nl
}

func (nl *NonLexical) Weight() float64 {
	return nl.weight
}

func (nl *NonLexical) Key() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(nl.Head))

	for _, b := range nl.Body {
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(b))
	}

	return sb.String()
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
