package pcfg

type lKey struct {
	i, n int
}

type rKey struct {
	n, j int
}

type Matcher struct {
	items map[Span]struct{}
	left  map[lKey][]*Item
	right map[rKey][]*Item
}

func NewMatcher() *Matcher {
	return &Matcher{
		items: make(map[Span]struct{}),
		left:  make(map[lKey][]*Item),
		right: make(map[rKey][]*Item),
	}
}

func (m *Matcher) Add(i *Item) bool {
	if _, ok := m.items[i.Span]; ok {
		return false
	} else {
		m.items[i.Span] = struct{}{}
	}

	lk := lKey{
		i: i.i,
		n: i.n,
	}

	rk := rKey{
		n: i.n,
		j: i.j,
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

func (m *Matcher) MatchLeft(i, n int) []*Item {
	lk := lKey{
		i: i,
		n: n,
	}

	items, ok := m.left[lk]

	if !ok {
		return []*Item{}
	}

	return items
}

func (m *Matcher) MatchRight(n, j int) []*Item {
	rk := rKey{
		n: n,
		j: j,
	}

	items, ok := m.right[rk]

	if !ok {
		return []*Item{}
	}

	return items
}
