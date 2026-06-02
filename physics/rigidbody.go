package physics

import "github.com/rhgrant10/sim/geo"

type RigidBody struct {
	Velocity    geo.Vec2
	Mass        float64
	Restitution float64
}

func (rb *RigidBody) InverseMass() float64 {
	if rb.Mass == 0.0 {
		return 0
	}
	return 1 / rb.Mass
}
