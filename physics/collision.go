package physics

import (
	"math"

	"github.com/rhgrant10/sim/geo"
)

type Collidable interface {
	GetCollider() Collider
	GetTransform() *geo.Transform
	GetRigidBody() *RigidBody
}

type CollisionInfo struct {
	Normal geo.Vec2
	Depth  float64
}

func epa(simplex Simplex, support func(geo.Vec2) geo.Vec2) CollisionInfo {
	const tolerance = 1e-6
	poly := NewPolytope(simplex)
	for {
		e := poly.ClosestEdge()
		p := support(e.normal)
		if p.Dot(e.normal)-e.dist < tolerance {
			return CollisionInfo{Normal: e.normal, Depth: e.dist}
		}
		poly.Insert(e.index, p)
	}
}

func gjk(a, b Collidable) (Simplex, bool) {
	ta := a.GetTransform()
	tb := b.GetTransform()
	ca := a.GetCollider()
	cb := b.GetCollider()

	// Minkowski difference support
	support := func(dir geo.Vec2) geo.Vec2 {
		return ca.Support(dir, *ta).Sub(cb.Support(dir.Negated(), *tb))
	}

	var s Simplex
	dir := geo.Vec2{X: 1, Y: 0} // any starting direction works
	s.Add(support(dir))
	dir = s.points[0].Negated()

	for {
		p := support(dir)
		if p.Dot(dir) < 0 {
			return s, false
		}
		s.Add(p)
		if s.Contains(&dir) {
			return s, true
		}
	}
}

func DetectCollision(a, b Collidable) (CollisionInfo, bool) {
	simplex, ok := gjk(a, b)
	if !ok {
		return CollisionInfo{}, false
	}
	ta, tb := a.GetTransform(), b.GetTransform()
	ca, cb := a.GetCollider(), b.GetCollider()
	support := func(dir geo.Vec2) geo.Vec2 {
		return ca.Support(dir, *ta).Sub(cb.Support(dir.Negated(), *tb))
	}
	info := epa(simplex, support)
	// EPA returns the outward normal of the Minkowski difference (A-B), which points
	// from a toward b. Negate so the returned normal points from b toward a,
	// matching the convention ResolveCollision expects.
	return CollisionInfo{Normal: info.Normal.Negated(), Depth: info.Depth}, true
}

func ResolveCollision(a, b Collidable, info CollisionInfo) {
	n := info.Normal
	inverseMassSum := a.GetRigidBody().InverseMass() + b.GetRigidBody().InverseMass()

	av, bv := a.GetRigidBody().Velocity, b.GetRigidBody().Velocity
	relVel := av.Sub(bv).Dot(n)

	if relVel < 0 {
		restitution := (a.GetRigidBody().Restitution + b.GetRigidBody().Restitution) / 2
		if inverseMassSum > 0 {
			j := -(1 + restitution) * relVel / inverseMassSum
			impulse := n.Scale(j)
			a.GetRigidBody().Velocity = av.Add(impulse.Scale(a.GetRigidBody().InverseMass()))
			b.GetRigidBody().Velocity = bv.Sub(impulse.Scale(b.GetRigidBody().InverseMass()))
		}
	}

	// Positional correction: push overlapping bodies apart to prevent sinking.
	if inverseMassSum > 0 {
		const percent = 0.4
		const slop = 0.01
		corrMag := math.Max(info.Depth-slop, 0) / inverseMassSum * percent
		corr := n.Scale(corrMag)
		a.GetTransform().Pos = a.GetTransform().Pos.Add(corr.Scale(a.GetRigidBody().InverseMass()))
		b.GetTransform().Pos = b.GetTransform().Pos.Sub(corr.Scale(b.GetRigidBody().InverseMass()))
	}
}
