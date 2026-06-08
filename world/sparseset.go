package world

import (
	"iter"
)

const notPresent = -1

type SparseSet[T any] struct {
	sparse      []int
	dense       []T
	denseOwners []Entity
}

func (s *SparseSet[T]) Lookup(id Entity) (component *T, ok bool) {
	if s.sparseShortage(id) > 0 {
		return
	}
	i := s.sparse[id]
	if 0 <= i && i < len(s.dense) && s.denseOwners[i] == id {
		component = &s.dense[i]
		ok = true
	}
	return
}

func (s *SparseSet[T]) All() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for i := range s.dense {
			if !yield(&s.dense[i]) {
				return
			}
		}
	}
}

func (s *SparseSet[T]) Insert(id Entity, component T) {
	if shortage := s.sparseShortage(id); shortage > 0 {
		addition := make([]int, shortage)
		for i := range addition {
			addition[i] = -1
		}
		s.sparse = append(s.sparse, addition...)
	}
	i := len(s.dense)
	s.dense = append(s.dense, component)
	s.denseOwners = append(s.denseOwners, id)
	s.sparse[id] = i
}

func (s *SparseSet[T]) Delete(id Entity) {
	shortage := s.sparseShortage(id)
	if shortage > 0 {
		return
	}
	i := s.sparse[id]
	if i >= 0 && i < len(s.dense) && s.denseOwners[i] == id {
		dLast := len(s.dense) - 1
		dLastOwner := s.denseOwners[dLast]

		s.dense[i] = s.dense[dLast]
		s.dense = s.dense[:dLast]

		s.denseOwners[i] = s.denseOwners[dLast]
		s.denseOwners = s.denseOwners[:dLast]

		s.sparse[dLastOwner] = i
		s.sparse[id] = notPresent
	}
}

func (s *SparseSet[T]) sparseShortage(id Entity) int64 {
	maxId := int64(len(s.sparse) - 1)
	return int64(id) - maxId
}
