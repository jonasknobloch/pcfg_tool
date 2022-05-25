package grammar

import (
	"errors"
	"fmt"
	"hash/fnv"
	"pcfg_tool/internal/utility"
	"strings"
)

type NonLexical struct {
	Head   int
	Body   []int
	key    uint64
	weight float64
}

func NewNonLexical(head string, body []string, symbols *SymbolTable) (*NonLexical, error) {
	if len(body) == 0 {
		return nil, errors.New("empty body")
	}

	nl := &NonLexical{
		weight: 1,
		Body:   make([]int, len(body)),
	}

	if h, err := symbols.Atoi(head); err != nil {
		return nil, err
	} else {
		nl.Head = h
	}

	for i, b := range body {
		v, err := symbols.Atoi(b)

		if err != nil {
			return nil, err
		}

		nl.Body[i] = v
	}

	h := fnv.New64()

	_, _ = h.Write([]byte(head))

	for _, b := range body {
		_, _ = h.Write([]byte(b))
	}

	nl.key = h.Sum64()

	return nl, nil
}

func (nl *NonLexical) Key() uint64 {
	return nl.key
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

	sb.WriteString(fmt.Sprintf(" %s\n", utility.FormatWeight(nl.weight)))

	return sb.String(), nil
}
