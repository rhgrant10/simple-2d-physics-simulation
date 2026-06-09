package systems

import (
	"github.com/rhgrant10/sim/collision"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
)

type ColliderAdapter struct {
	transform *geo.Transform
	collider  *collision.Collider
	rigidbody *physics.RigidBody
	material  *physics.Material
}

func (a *ColliderAdapter) GetCollider() collision.Collider {
	return *a.collider
}

func (a *ColliderAdapter) GetTransform() *geo.Transform {
	return a.transform
}

func (a *ColliderAdapter) GetRigidBody() *physics.RigidBody {
	return a.rigidbody
}

func (a *ColliderAdapter) GetMaterial() *physics.Material {
	return a.material
}
