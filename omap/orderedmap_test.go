package omap

import "testing"

func TestOrderedMapBasic(t *testing.T) {
	m := New[string, int]()

	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	if m.Len() != 3 {
		t.Errorf("Expected len 3, got %d", m.Len())
	}

	if opt := m.Get("b"); !opt.IsPresent() || opt.MustGet() != 2 {
		t.Errorf("Expected b=2, got %v", opt)
	}

	if !m.Has("a") {
		t.Error("Expected to have key 'a'")
	}
}

func TestOrderedMapOrder(t *testing.T) {
	m := New[string, int]()
	m.Set("z", 1)
	m.Set("a", 2)
	m.Set("m", 3)

	keys := m.Keys()
	expected := []string{"z", "a", "m"}

	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Expected key %s at %d, got %s", expected[i], i, key)
		}
	}
}

func TestOrderedMapUpdateMovesToEnd(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.Set("a", 10)

	keys := m.Keys()
	expected := []string{"b", "c", "a"}

	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Expected key %s at %d, got %s", expected[i], i, key)
		}
	}

	if opt := m.Get("a"); opt.MustGet() != 10 {
		t.Errorf("Expected a=10, got %v", opt.Get())
	}
}

func TestOrderedMapFirstLast(t *testing.T) {
	m := New[string, int]()
	m.Set("first", 1)
	m.Set("middle", 2)
	m.Set("last", 3)

	optFirst := m.First()
	if !optFirst.IsPresent() || optFirst.MustGet().Key != "first" || optFirst.MustGet().Value != 1 {
		t.Errorf("Expected first=1, got %v", optFirst)
	}

	optLast := m.Last()
	if !optLast.IsPresent() || optLast.MustGet().Key != "last" || optLast.MustGet().Value != 3 {
		t.Errorf("Expected last=3, got %v", optLast)
	}
}

func TestOrderedMapPopFirstLast(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	optFirst := m.PopFirst()
	if !optFirst.IsPresent() || optFirst.MustGet().Key != "a" || optFirst.MustGet().Value != 1 {
		t.Errorf("PopFirst failed: %v", optFirst)
	}

	optLast := m.PopLast()
	if !optLast.IsPresent() || optLast.MustGet().Key != "c" || optLast.MustGet().Value != 3 {
		t.Errorf("PopLast failed: %v", optLast)
	}

	if m.Len() != 1 {
		t.Errorf("Expected len 1, got %d", m.Len())
	}
}

func TestOrderedMapDelete(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	if !m.Delete("b") {
		t.Error("Delete should return true")
	}

	if m.Has("b") {
		t.Error("Should not have key 'b' after delete")
	}

	keys := m.Keys()
	expected := []string{"a", "c"}

	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Expected key %s at %d, got %s", expected[i], i, key)
		}
	}
}

func TestOrderedMapRange(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	count := 0
	m.Range(func(k string, v int) bool {
		count++
		return true
	})

	if count != 3 {
		t.Errorf("Expected 3 iterations, got %d", count)
	}
}
