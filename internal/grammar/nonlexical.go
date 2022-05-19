package grammar

import "strings"

type NonLexical struct {
	Head string
	Body []string
}

func (nl *NonLexical) KeyHead() string {
	return nl.Head
}

func (nl *NonLexical) KeyBody() string {
	return strings.Join(nl.Body, " ")
}

func (nl *NonLexical) String() string {
	return nl.KeyHead() + " -> " + nl.KeyBody()
}
