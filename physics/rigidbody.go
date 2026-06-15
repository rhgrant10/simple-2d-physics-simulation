package physics

import "github.com/rhgrant10/sim/geo"

type RigidBody struct {
	Velocity geo.Vec2
	PrevPos  geo.Vec2
	Mass     float64
}

func (rb *RigidBody) InverseMass() float64 {
	if rb.IsStatic() {
		return 0
	}
	return 1 / rb.Mass
}

func (rb *RigidBody) IsStatic() bool {
	return rb.Mass == 0.0
}
