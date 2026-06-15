package world

import (
	"iter"
)

type View2[A, B any] struct {
	Entity Entity
	A      *A
	B      *B
}

type View3[A, B, C any] struct {
	Entity Entity
	A      *A
	B      *B
	C      *C
}

type View4[A, B, C, D any] struct {
	Entity Entity
	A      *A
	B      *B
	C      *C
	D      *D
}

func Query2[A, B any](w *World) iter.Seq[*View2[A, B]] {
	setA := getSet[A](w)
	setB := getSet[B](w)
	return func(yield func(*View2[A, B]) bool) {
		for each := range intersection2(setA, setB) {
			if !yield(each) {
				return
			}
		}
	}
}

func Query3[A, B, C any](w *World) iter.Seq[*View3[A, B, C]] {
	setA := getSet[A](w)
	setB := getSet[B](w)
	setC := getSet[C](w)
	return func(yield func(*View3[A, B, C]) bool) {
		for each := range intersection3(setA, setB, setC) {
			if !yield(each) {
				return
			}
		}
	}
}

func Query4[A, B, C, D any](w *World) iter.Seq[*View4[A, B, C, D]] {
	setA := getSet[A](w)
	setB := getSet[B](w)
	setC := getSet[C](w)
	setD := getSet[D](w)
	return func(yield func(*View4[A, B, C, D]) bool) {
		for each := range intersection4(setA, setB, setC, setD) {
			if !yield(each) {
				return
			}
		}
	}
}

func HasComponent[T any](w *World, id Entity) bool {
	_, ok := getSet[T](w).Lookup(id)
	return ok
}

func GetComponent[T any](w *World, id Entity) (component *T, ok bool) {
	return getSet[T](w).Lookup(id)
}

func intersection2[A, B any](a *SparseSet[A], b *SparseSet[B]) iter.Seq[*View2[A, B]] {
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

func intersection3[A, B, C any](a *SparseSet[A], b *SparseSet[B], c *SparseSet[C]) iter.Seq[*View3[A, B, C]] {
	return func(yield func(*View3[A, B, C]) bool) {
		for i := range a.denseOwners {
			e := a.denseOwners[i]
			bComponent, bOk := b.Lookup(e)
			if !bOk {
				continue
			}
			cComponent, cOk := c.Lookup(e)
			if !cOk {
				continue
			}
			v := View3[A, B, C]{Entity: e, A: &a.dense[i], B: bComponent, C: cComponent}
			if !yield(&v) {
				return
			}
		}
	}
}

func intersection4[A, B, C, D any](a *SparseSet[A], b *SparseSet[B], c *SparseSet[C], d *SparseSet[D]) iter.Seq[*View4[A, B, C, D]] {
	return func(yield func(*View4[A, B, C, D]) bool) {
		for i := range a.denseOwners {
			e := a.denseOwners[i]
			bComponent, bOk := b.Lookup(e)
			if !bOk {
				continue
			}
			cComponent, cOk := c.Lookup(e)
			if !cOk {
				continue
			}
			dComponent, dOk := d.Lookup(e)
			if !dOk {
				continue
			}
			v := View4[A, B, C, D]{Entity: e, A: &a.dense[i], B: bComponent, C: cComponent, D: dComponent}
			if !yield(&v) {
				return
			}
		}
	}
}
