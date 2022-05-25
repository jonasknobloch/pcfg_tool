package grammar

import (
	"fmt"
	"hash/fnv"
	"pcfg_tool/internal/utility"
	"strings"
)

type Lexical struct {
	Head   int
	Body   string
	key    uint64
	weight float64
}

func NewLexical(head, body string, symbols *SymbolTable) (*Lexical, error) {
	l := &Lexical{
		weight: 1,
		Body:   body,
	}

	if h, err := symbols.Atoi(head); err != nil {
		return nil, err
	} else {
		l.Head = h
	}

	h := fnv.New64()

	_, _ = h.Write([]byte(head))
	_, _ = h.Write([]byte(body))

	l.key = h.Sum64()

	return l, nil
}

func (l *Lexical) Key() uint64 {
	return l.key
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

	sb.WriteString(fmt.Sprintf(" %s\n", utility.FormatWeight(l.weight)))

	return sb.String(), nil
}
