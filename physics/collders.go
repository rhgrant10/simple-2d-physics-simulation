package physics

import (
	"math"

	"github.com/rhgrant10/sim/geo"
)

type Collider interface {
	Support(dir geo.Vec2, transform geo.Transform) geo.Vec2
}

type CircleCollider struct {
	Radius float64
}

func (c *CircleCollider) Support(dir geo.Vec2, transform geo.Transform) geo.Vec2 {
	n := dir.Normalized()
	return geo.Vec2{
		X: transform.Pos.X + n.X*c.Radius,
		Y: transform.Pos.Y + n.Y*c.Radius,
	}
}

type BoxCollider struct {
	HalfW float64
	HalfH float64
}

func (b *BoxCollider) Support(dir geo.Vec2, transform geo.Transform) geo.Vec2 {
	return geo.Vec2{
		X: transform.Pos.X + math.Copysign(b.HalfW, dir.X),
		Y: transform.Pos.Y + math.Copysign(b.HalfH, dir.Y),
	}
}
