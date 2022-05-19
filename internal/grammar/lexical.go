package grammar

type Lexical struct {
	Head, Body string
}

func (l *Lexical) KeyHead() string {
	return l.Head
}

func (l *Lexical) KeyBody() string {
	return l.Body
}

func (l *Lexical) String() string {
	return l.KeyHead() + " " + l.KeyBody()
}
