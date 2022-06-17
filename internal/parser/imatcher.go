package parser

type ItemMatcher struct {
	items map[Span]struct{}
	left  map[int][]*Item
	right map[int][]*Item
}

func NewItemMatcher() *ItemMatcher {
	return &ItemMatcher{
		items: make(map[Span]struct{}),
		left:  make(map[int][]*Item),
		right: make(map[int][]*Item),
	}
}

func (m *ItemMatcher) Add(i *Item) bool {
	if _, ok := m.items[i.Span]; ok {
		return false
	} else {
		m.items[i.Span] = struct{}{}
	}

	if _, ok := m.left[i.i]; !ok {
		m.left[i.i] = make([]*Item, 0)
	}

	if _, ok := m.right[i.j]; !ok {
		m.right[i.j] = make([]*Item, 0)
	}

	m.left[i.i] = append(m.left[i.i], i)
	m.right[i.j] = append(m.right[i.j], i)

	return true
}

func (m *ItemMatcher) Match(item *Item) ([]*Item, []*Item) {
	return m.right[item.i], m.left[item.j]
}
