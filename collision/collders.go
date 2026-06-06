package collision

import (
	"math"

	"github.com/rhgrant10/sim/geo"
)

type Collider interface {
	Support(dir geo.Vec2, transform geo.Transform) geo.Vec2
	AABB(transform geo.Transform) geo.AABB
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

func (c *CircleCollider) AABB(transform geo.Transform) geo.AABB {
	squareRadius := geo.Vec2{X: c.Radius, Y: c.Radius}
	return geo.NewAABB(transform.Pos.Sub(squareRadius), transform.Pos.Add(squareRadius))
}

// Assumes no rotation
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

func (b *BoxCollider) AABB(transform geo.Transform) geo.AABB {
	halfSize := geo.Vec2{X: b.HalfW, Y: b.HalfH}
	return geo.NewAABB(transform.Pos.Sub(halfSize), transform.Pos.Add(halfSize))
}
