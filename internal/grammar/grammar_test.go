package grammar

import "testing"

type rule struct{}

func (r *rule) Weight() float64 {
	return 0
}

func (r *rule) Key() string {
	return ""
}

func (r *rule) String(st *SymbolTable) (string, error) {
	return "", nil
}

func TestGrammar_Initial(t *testing.T) {
	g := NewGrammar()

	if foo, bar := g.initial, g.Initial(); foo != bar {
		t.Errorf("expected %d but got %d", foo, bar)
	}
}

func TestGrammar_SetInitial(t *testing.T) {
	g := NewGrammar()

	foo := "FOO"

	g.SetInitial(foo)

	bar, err := g.Symbols.Itoa(g.Initial())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if foo != bar {
		t.Errorf("expected %s but got %s", foo, bar)
	}
}

func TestGrammar_AddRule(t *testing.T) {
	g := NewGrammar()

	if err := g.AddRule(NewNonLexical("FOO", []string{"foo"}, 1, g.Symbols)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := g.AddRule(NewLexical("BAR", "bar", 1, g.Symbols)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := g.AddRule(&rule{}); err == nil {
		t.Fatalf("add rule should have errored")
	}
}

func TestGrammar_AddNonLexical(t *testing.T) {
	g := NewGrammar()

	g.AddNonLexical(NewNonLexical("FOO", []string{"foo"}, 0.5, g.Symbols))
	g.AddNonLexical(NewNonLexical("FOO", []string{"foo"}, 0.5, g.Symbols))

	rules := g.Rules(g.Symbols.Atoi("foo"))

	if len(rules) != 1 {
		t.Fatalf("unexpected lexion length")
	}

	if rules[0].weight != 1 {
		t.Fatalf("unexpeted lexical weight")
	}
}

func TestGrammar_AddLexical(t *testing.T) {
	g := NewGrammar()

	g.AddLexical(NewLexical("FOO", "foo", 0.5, g.Symbols))
	g.AddLexical(NewLexical("FOO", "foo", 0.5, g.Symbols))

	lexicon := g.Lexicon("foo")

	if len(lexicon) != 1 {
		t.Fatalf("unexpected lexion length")
	}

	if lexicon[0].weight != 1 {
		t.Fatalf("unexpeted lexical weight")
	}
}

func TestGrammar_Normalize(t *testing.T) {
	g := NewGrammar()

	g.AddNonLexical(NewNonLexical("ROOT", []string{"FOO", "BAR"}, 1, g.Symbols))

	g.AddLexical(NewLexical("FOO", "foo", 1, g.Symbols))
	g.AddLexical(NewLexical("BAR", "bar", 1, g.Symbols))
	g.AddLexical(NewLexical("BAR", "baz", 1, g.Symbols))

	if g.IsNormalized() {
		t.Fatalf("not normalized grammar classified as normalized")
	}

	g.Normalize()

	if !g.IsNormalized() {
		t.Fatalf("normalized grammar classified as not normalized")
	}
}

func TestGrammar_IsNormalized(t *testing.T) {
	g := NewGrammar()

	g.AddNonLexical(NewNonLexical("ROOT", []string{"FOO", "BAR"}, 1, g.Symbols))

	g.AddLexical(NewLexical("FOO", "foo", 1, g.Symbols))
	g.AddLexical(NewLexical("BAR", "bar", 0.5, g.Symbols))

	if g.IsNormalized() {
		t.Fatalf("not normalized grammar classified as normalized")
	}

	g.AddLexical(NewLexical("BAR", "baz", 0.5, g.Symbols))

	if !g.IsNormalized() {
		t.Fatalf("normalized grammar classified as not normalized")
	}
}

func TestGrammar_Rules(t *testing.T) {
	g := NewGrammar()

	rules := []*NonLexical{
		NewNonLexical("ROOT", []string{"FOO", "BAR"}, 1, g.Symbols),
		NewNonLexical("ROOT", []string{"FOO", "BAZ"}, 1, g.Symbols),
	}

	for _, r := range rules {
		g.AddNonLexical(r)
	}

	if len(g.Rules(0)) != 0 {
		t.Fatalf("rules should be empty")
	}

	for i, r := range g.Rules(g.Symbols.Atoi("FOO")) {
		if r != rules[i] {
			t.Fatalf("expected %v but got %v", rules[i], r)
		}
	}
}

func TestGrammar_Lexicon(t *testing.T) {
	g := NewGrammar()

	rules := []*Lexical{
		NewLexical("FOO", "foo", 1, g.Symbols),
		NewLexical("BAR", "foo", 1, g.Symbols),
	}

	for _, r := range rules {
		g.AddLexical(r)
	}

	if len(g.Lexicon("")) != 0 {
		t.Fatalf("lexicon should be empty")
	}

	for i, r := range g.Lexicon("foo") {
		if r != rules[i] {
			t.Fatalf("expected %v but got %v", rules[i], r)
		}
	}
}

func TestGrammar_Import(t *testing.T) {
	// TODO implement
}

func TestGrammar_Export(t *testing.T) {
	// TODO implement
}
