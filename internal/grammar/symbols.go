package grammar

import "errors"

type SymbolTable struct {
	index int
	atoi  map[string]int
	itoa  map[int]string
}

var ErrUnknownSymbol = errors.New("unknown symbol")

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		index: 0,
		atoi:  make(map[string]int),
		itoa:  make(map[int]string),
	}
}

func (st *SymbolTable) Atoi(s string) int {
	if i, ok := st.atoi[s]; ok {
		return i
	}

	st.index++

	st.itoa[st.index] = s
	st.atoi[s] = st.index

	return st.index
}

func (st *SymbolTable) Itoa(i int) (string, error) {
	a, ok := st.itoa[i]

	if !ok {
		return "", ErrUnknownSymbol
	}

	return a, nil
}
