package collision

import "github.com/rhgrant10/sim/geo"

type Detector interface {
	Detect(entities []Collidable) []CollisionInfo
}

type DefaultDetector struct {
	broad      Broadphaser
	narrow     NarrowPhaser
	closePairs []geo.PairKey
}

func NewDefaultDetector(broad Broadphaser, narrow NarrowPhaser) DefaultDetector {
	return DefaultDetector{broad: broad, narrow: narrow}
}

func (h *DefaultDetector) Detect(entities []Collidable) []CollisionInfo {
	staticEntities := []Collidable{}
	dynamicEntities := []Collidable{}
	for i := range entities {
		if entities[i].GetRigidBody().IsStatic() {
			staticEntities = append(staticEntities, entities[i])
		} else {
			dynamicEntities = append(dynamicEntities, entities[i])
		}
	}

	// collisions with dynamic objects - build
	bvs := make([]geo.AABB, 0, len(dynamicEntities))
	for i := range dynamicEntities {
		if !dynamicEntities[i].GetRigidBody().IsStatic() {
			transform := *dynamicEntities[i].GetTransform()
			aabb := dynamicEntities[i].GetCollider().AABB(transform)
			bvs = append(bvs, aabb)
		}
	}

	h.broad.Clear()
	h.broad.Build(&bvs)

	// collisions with dynamic objects - find & process
	collisions := []CollisionInfo{}
	h.closePairs = h.broad.CandidatePairs(h.closePairs)
	for _, p := range h.closePairs {
		a, b := p.A(), p.B()
		if info, ok := h.narrow.Detect(dynamicEntities[a], dynamicEntities[b]); ok {
			collisions = append(collisions, info)
		}
	}

	// collisions with static objects
	for d := range dynamicEntities {
		for s := range staticEntities {
			if info, ok := h.narrow.Detect(dynamicEntities[d], staticEntities[s]); ok {
				collisions = append(collisions, info)
			}
		}
	}
	return collisions
}
