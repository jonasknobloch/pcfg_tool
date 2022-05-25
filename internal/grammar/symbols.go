package grammar

import (
	"errors"
	"sync"
)

type SymbolTable struct {
	index int
	atoi  map[string]int
	itoa  map[int]string
	mutex sync.RWMutex
}

var ErrUnknownSymbol = errors.New("unknown symbol")

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		index: 0,
		atoi:  make(map[string]int),
		itoa:  make(map[int]string),
		mutex: sync.RWMutex{},
	}
}

func (st *SymbolTable) Atoi(s string) (int, error) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if i, ok := st.atoi[s]; ok {
		return i, nil
	}

	st.index++

	st.itoa[st.index] = s
	st.atoi[s] = st.index

	return st.index, nil
}

func (st *SymbolTable) Itoa(i int) (string, error) {
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	a, ok := st.itoa[i]

	if !ok {
		return "", ErrUnknownSymbol
	}

	return a, nil
}
