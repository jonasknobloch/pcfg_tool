package grammar

import (
	"pcfg_tool/internal/utility"
	"strings"
)

type Lexical struct {
	Head   int
	Body   string
	weight float64
}

func NewLexical(head, body string, symbols *SymbolTable) (*Lexical, string, error) {
	l := &Lexical{
		weight: 1,
		Body:   body,
	}

	if h, err := symbols.Atoi(head); err != nil {
		return nil, "", err
	} else {
		l.Head = h
	}

	key := head + " " + body

	return l, key, nil
}

func (l *Lexical) Weight() float64 {
	return l.weight
}

func (l *Lexical) String(st *SymbolTable) (string, error) {
	var sb strings.Builder

	if head, err := st.Itoa(l.Head); err != nil {
		return "", err
	} else {
		sb.WriteString(head)
	}

	sb.WriteString(" ")
	sb.WriteString(l.Body)

	sb.WriteString(" ")
	sb.WriteString(utility.FormatWeight(l.weight))

	return sb.String(), nil
}
