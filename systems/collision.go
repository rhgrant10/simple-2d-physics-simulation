package systems

import (
	"github.com/rhgrant10/sim/collision"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
	"github.com/rhgrant10/sim/world"
)

type CollisionInfo struct {
	A      world.Entity
	B      world.Entity
	Normal geo.Vec2
	Depth  float64
}

type BasicCollisionSystem struct {
	Detector collision.Detector
	Resolver collision.Resolver
}

func (s *BasicCollisionSystem) Solve(w *world.World) {
	entities := []collision.Collidable{}
	for c := range world.Query3[collision.Collider, geo.Transform, physics.Material](w) {
		rb, _ := world.GetComponent[physics.RigidBody](w, c.Entity)
		entities = append(entities, &ColliderAdapter{
			collider:  c.A,
			transform: c.B,
			material:  c.C,
			rigidbody: rb,
		})
	}
	for _, c := range s.Detector.Detect(entities) {
		s.Resolver.Resolve(c)
	}
}
