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

	if val, ok := m.Get("b"); !ok || val != 2 {
		t.Errorf("Expected b=2, got %v, %v", val, ok)
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

	if val, _ := m.Get("a"); val != 10 {
		t.Errorf("Expected a=10, got %d", val)
	}
}

func TestOrderedMapFirstLast(t *testing.T) {
	m := New[string, int]()
	m.Set("first", 1)
	m.Set("middle", 2)
	m.Set("last", 3)

	k, v, ok := m.First()
	if !ok || k != "first" || v != 1 {
		t.Errorf("Expected first=1, got %s=%d, %v", k, v, ok)
	}

	k, v, ok = m.Last()
	if !ok || k != "last" || v != 3 {
		t.Errorf("Expected last=3, got %s=%d, %v", k, v, ok)
	}
}

func TestOrderedMapPopFirstLast(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	k, v, ok := m.PopFirst()
	if !ok || k != "a" || v != 1 {
		t.Errorf("PopFirst failed: %s=%d, %v", k, v, ok)
	}

	k, v, ok = m.PopLast()
	if !ok || k != "c" || v != 3 {
		t.Errorf("PopLast failed: %s=%d, %v", k, v, ok)
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
