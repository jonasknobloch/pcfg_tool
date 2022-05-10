package pcfg

type lKey struct {
	int
	string
}

type rKey struct {
	string
	int
}

type Matcher struct {
	items map[Item]struct{}
	left  map[lKey][]*Item
	right map[rKey][]*Item
}

func NewMatcher() *Matcher {
	return &Matcher{
		items: make(map[Item]struct{}),
		left:  make(map[lKey][]*Item),
		right: make(map[rKey][]*Item),
	}
}

func (m *Matcher) Add(i *Item) bool {
	if _, ok := m.items[*i]; ok {
		return false
	} else {
		m.items[*i] = struct{}{}
	}

	lk := lKey{
		int:    i.i,
		string: i.n,
	}

	rk := rKey{
		string: i.n,
		int:    i.j,
	}

	if _, ok := m.left[lk]; !ok {
		m.left[lk] = make([]*Item, 0)
	}

	if _, ok := m.right[rk]; !ok {
		m.right[rk] = make([]*Item, 0)
	}

	m.left[lk] = append(m.left[lk], i)
	m.right[rk] = append(m.right[rk], i)

	return true
}

func (m *Matcher) MatchLeft(i int, n string) []*Item {
	lk := lKey{
		int:    i,
		string: n,
	}

	items, ok := m.left[lk]

	if !ok {
		return []*Item{}
	}

	return items
}

func (m *Matcher) MatchRight(n string, j int) []*Item {
	rk := rKey{
		string: n,
		int:    j,
	}

	items, ok := m.right[rk]

	if !ok {
		return []*Item{}
	}

	return items
}
