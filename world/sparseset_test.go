package world

import (
	"slices"
	"testing"
)

// collect drains All() into a slice of values (not pointers) for easy comparison.
func collect[T any](s *SparseSet[T]) []T {
	var out []T
	for v := range s.All() {
		out = append(out, *v)
	}
	return out
}

func TestSparseSet_LookupMissing(t *testing.T) {
	var s SparseSet[int]
	if _, ok := s.Lookup(0); ok {
		t.Error("lookup on empty set should return false")
	}
	if _, ok := s.Lookup(999); ok {
		t.Error("lookup of never-inserted entity should return false")
	}
}

func TestSparseSet_InsertThenLookup(t *testing.T) {
	var s SparseSet[int]
	s.Insert(5, 42)
	v, ok := s.Lookup(5)
	if !ok {
		t.Fatal("expected to find entity 5")
	}
	if *v != 42 {
		t.Errorf("got %d, want 42", *v)
	}
}

func TestSparseSet_LookupReturnsMutablePointer(t *testing.T) {
	var s SparseSet[int]
	s.Insert(0, 10)
	v, _ := s.Lookup(0)
	*v = 99
	v2, _ := s.Lookup(0)
	if *v2 != 99 {
		t.Error("mutation via pointer should be reflected in subsequent lookup")
	}
}

func TestSparseSet_DeleteThenLookup(t *testing.T) {
	var s SparseSet[int]
	s.Insert(3, 7)
	s.Delete(3)
	if _, ok := s.Lookup(3); ok {
		t.Error("lookup after delete should return false")
	}
}

func TestSparseSet_DeleteLastElement(t *testing.T) {
	var s SparseSet[int]
	s.Insert(0, 1)
	s.Delete(0)
	if _, ok := s.Lookup(0); ok {
		t.Error("lookup of sole deleted element should return false")
	}
	if len(collect(&s)) != 0 {
		t.Error("All() should be empty after deleting sole element")
	}
}

// Deleting entity A (not the last) swaps the last entity into A's slot.
// The swapped entity must still be findable afterwards.
func TestSparseSet_DeletePreservesSwappedEntity(t *testing.T) {
	var s SparseSet[int]
	s.Insert(0, 10)
	s.Insert(1, 20)
	s.Insert(2, 30)

	s.Delete(0) // entity 2 (last) gets swapped into slot 0

	if _, ok := s.Lookup(0); ok {
		t.Error("deleted entity 0 should not be found")
	}
	v, ok := s.Lookup(2)
	if !ok {
		t.Fatal("swapped entity 2 should still be found")
	}
	if *v != 30 {
		t.Errorf("swapped entity 2 has wrong value: got %d, want 30", *v)
	}
	v, ok = s.Lookup(1)
	if !ok {
		t.Fatal("untouched entity 1 should still be found")
	}
	if *v != 20 {
		t.Errorf("entity 1 has wrong value: got %d, want 20", *v)
	}
}

func TestSparseSet_DeleteNonexistent(t *testing.T) {
	var s SparseSet[int]
	s.Insert(0, 1)
	s.Delete(999) // should not panic or corrupt
	if _, ok := s.Lookup(0); !ok {
		t.Error("unrelated entity should be unaffected by delete of nonexistent")
	}
}

func TestSparseSet_AllYieldsAllValues(t *testing.T) {
	var s SparseSet[int]
	s.Insert(0, 10)
	s.Insert(5, 20)
	s.Insert(100, 30)

	got := collect(&s)
	slices.Sort(got)
	want := []int{10, 20, 30}
	if !slices.Equal(got, want) {
		t.Errorf("All() yielded %v, want %v", got, want)
	}
}

func TestSparseSet_AllAfterDelete(t *testing.T) {
	var s SparseSet[int]
	s.Insert(0, 10)
	s.Insert(1, 20)
	s.Insert(2, 30)
	s.Delete(1)

	got := collect(&s)
	slices.Sort(got)
	want := []int{10, 30}
	if !slices.Equal(got, want) {
		t.Errorf("All() after delete yielded %v, want %v", got, want)
	}
}

func TestSparseSet_LargeEntityId(t *testing.T) {
	var s SparseSet[int]
	s.Insert(1000, 42)
	v, ok := s.Lookup(1000)
	if !ok {
		t.Fatal("entity with large ID should be found")
	}
	if *v != 42 {
		t.Errorf("got %d, want 42", *v)
	}
	if _, ok := s.Lookup(500); ok {
		t.Error("entity 500 was never inserted, should not be found")
	}
}
