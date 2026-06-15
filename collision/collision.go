package collision

import (
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
)

type Collidable interface {
	GetCollider() Collider
	GetTransform() *geo.Transform
	GetRigidBody() *physics.RigidBody
	GetMaterial() *physics.Material
}

type CollisionInfo struct {
	A                Collidable
	B                Collidable
	Normal           geo.Vec2
	Depth            float64
	RelativeVelocity geo.Vec2
}
