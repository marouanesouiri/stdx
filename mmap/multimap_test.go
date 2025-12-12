package mmap

import "testing"

func TestMultimapBasic(t *testing.T) {
	m := New[string, int]()

	added := m.Put("key", 1)
	if !added {
		t.Error("Expected Put to return true for new value")
	}

	added = m.Put("key", 1)
	if added {
		t.Error("Expected Put to return false for duplicate value")
	}

	added = m.Put("key", 2)
	if !added {
		t.Error("Expected Put to return true for new value")
	}

	if m.Size() != 2 {
		t.Errorf("Expected size 2, got %d", m.Size())
	}

	if m.Len() != 1 {
		t.Errorf("Expected 1 key, got %d", m.Len())
	}
}

func TestMultimapGet(t *testing.T) {
	m := New[string, string]()
	m.Put("tags", "go")
	m.Put("tags", "generics")
	m.Put("tags", "stdlib")

	values := m.Get("tags")
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	empty := m.Get("nonexistent")
	if len(empty) != 0 {
		t.Errorf("Expected empty slice, got %v", empty)
	}
}

func TestMultimapPutAll(t *testing.T) {
	m := New[string, int]()
	count := m.PutAll("nums", 1, 2, 3, 1)

	if count != 3 {
		t.Errorf("Expected 3 new values, got %d", count)
	}

	if m.KeySize("nums") != 3 {
		t.Errorf("Expected 3 values for key, got %d", m.KeySize("nums"))
	}
}

func TestMultimapDelete(t *testing.T) {
	m := New[string, int]()
	m.Put("nums", 1)
	m.Put("nums", 2)
	m.Put("nums", 3)

	deleted := m.Delete("nums", 2)
	if !deleted {
		t.Error("Expected Delete to return true")
	}

	if m.Contains("nums", 2) {
		t.Error("Should not contain deleted value")
	}

	if m.KeySize("nums") != 2 {
		t.Errorf("Expected 2 values, got %d", m.KeySize("nums"))
	}
}

func TestMultimapDeleteAll(t *testing.T) {
	m := New[string, int]()
	m.Put("nums", 1)
	m.Put("nums", 2)
	m.Put("nums", 3)

	deleted := m.DeleteAll("nums")
	if !deleted {
		t.Error("Expected DeleteAll to return true")
	}

	if m.ContainsKey("nums") {
		t.Error("Should not contain key after DeleteAll")
	}

	if m.Size() != 0 {
		t.Errorf("Expected size 0, got %d", m.Size())
	}
}

func TestMultimapContains(t *testing.T) {
	m := New[string, string]()
	m.Put("color", "red")
	m.Put("color", "blue")

	if !m.ContainsKey("color") {
		t.Error("Should contain key 'color'")
	}

	if !m.Contains("color", "red") {
		t.Error("Should contain entry color=red")
	}

	if m.Contains("color", "green") {
		t.Error("Should not contain entry color=green")
	}
}

func TestMultimapRange(t *testing.T) {
	m := New[string, int]()
	m.Put("a", 1)
	m.Put("a", 2)
	m.Put("b", 3)

	count := 0
	m.Range(func(k string, v int) bool {
		count++
		return true
	})

	if count != 3 {
		t.Errorf("Expected 3 iterations, got %d", count)
	}
}

func TestMultimapForEachKey(t *testing.T) {
	m := New[string, int]()
	m.Put("a", 1)
	m.Put("a", 2)
	m.Put("b", 3)

	keyCount := 0
	m.ForEachKey(func(k string, vals []int) bool {
		keyCount++
		if k == "a" && len(vals) != 2 {
			t.Errorf("Expected 2 values for key 'a', got %d", len(vals))
		}
		return true
	})

	if keyCount != 2 {
		t.Errorf("Expected 2 keys, got %d", keyCount)
	}
}
