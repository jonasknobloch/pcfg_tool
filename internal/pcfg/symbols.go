package pcfg

import (
	"hash/fnv"
)

type SymbolTable struct {
	atoi map[string]int
	itoa map[int]string
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		atoi: make(map[string]int),
		itoa: make(map[int]string),
	}
}

func (st *SymbolTable) Atoi(s string) (int, error) {
	if i, ok := st.atoi[s]; ok {
		return i, nil
	}

	h := fnv.New32a()

	if _, err := h.Write([]byte(s)); err != nil {
		return 0, err
	}

	i := int(h.Sum32())

	st.itoa[i] = s
	st.atoi[s] = i

	return i, nil
}

func (st *SymbolTable) Itoa(i int) (string, bool) {
	a, ok := st.itoa[i]

	return a, ok
}
