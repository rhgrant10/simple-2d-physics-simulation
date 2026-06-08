package world

import (
	"iter"
)

type View2[A, B any] struct {
	Entity Entity
	A      *A
	B      *B
}

func Query2[A, B any](w *World) iter.Seq[*View2[A, B]] {
	setA := getSet[A](w)
	setB := getSet[B](w)
	return func(yield func(*View2[A, B]) bool) {
		for each := range intersection(setA, setB) {
			if !yield(each) {
				return
			}
		}
	}
}

func intersection[A, B any](a *SparseSet[A], b *SparseSet[B]) iter.Seq[*View2[A, B]] {
	return func(yield func(*View2[A, B]) bool) {
		for i := range a.denseOwners {
			e := a.denseOwners[i]
			if bComponent, ok := b.Lookup(e); ok {
				v := View2[A, B]{Entity: e, A: &a.dense[i], B: bComponent}
				if !yield(&v) {
					return
				}
			}
		}
	}
}
