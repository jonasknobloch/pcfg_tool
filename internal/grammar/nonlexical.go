package grammar

import (
	"errors"
	"pcfg_tool/internal/utility"
	"strings"
)

type NonLexical struct {
	Head   int
	Body   []int
	weight float64
}

func NewNonLexical(head string, body []string, symbols *SymbolTable) (*NonLexical, string, error) {
	if len(body) == 0 {
		return nil, "", errors.New("empty body")
	}

	nl := &NonLexical{
		weight: 1,
		Body:   make([]int, len(body)),
	}

	if h, err := symbols.Atoi(head); err != nil {
		return nil, "", err
	} else {
		nl.Head = h
	}

	for i, b := range body {
		v, err := symbols.Atoi(b)

		if err != nil {
			return nil, "", err
		}

		nl.Body[i] = v
	}

	key := head + " " + strings.Join(body, " ")

	return nl, key, nil
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
