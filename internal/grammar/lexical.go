package grammar

import (
	"pcfg_tool/internal/utility"
	"strconv"
	"strings"
)

type Lexical struct {
	Head   int
	Body   string
	weight float64
}

func NewLexical(head, body string, weight float64, symbols *SymbolTable) *Lexical {
	l := &Lexical{
		weight: weight,
		Body:   body,
	}

	l.Head = symbols.Atoi(head)

	return l
}

func (l *Lexical) Weight() float64 {
	return l.weight
}

func (l *Lexical) Key() string {
	return strconv.Itoa(l.Head) + " " + l.Body
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
