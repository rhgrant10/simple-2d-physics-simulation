package world

import (
	"reflect"
)

type Entity uint64

type World struct {
	components   map[reflect.Type]any
	nextEntityId Entity
}

func NewWorld() *World {
	return &World{
		components:   make(map[reflect.Type]any),
		nextEntityId: 0,
	}
}

func Spawn(w *World) Entity {
	entity := w.nextEntityId
	w.nextEntityId++
	return entity
}

func AddComponent[T any](w *World, id Entity, c T) {
	set := getSet[T](w)
	set.Insert(id, c)
}

func getSet[T any](w *World) *SparseSet[T] {
	t := reflect.TypeFor[*T]()
	set, ok := w.components[t]
	if !ok {
		set = &SparseSet[T]{}
		w.components[t] = set
	}
	return set.(*SparseSet[T])
}
