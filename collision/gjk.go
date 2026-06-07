package collision

import "github.com/rhgrant10/sim/geo"

const epaTolerance = 1e-6

type GJK struct {
}

func (d *GJK) Check(a, b Collidable) (info CollisionInfo, areColliding bool) {
	simplex, areColliding := gjk(a, b)
	if areColliding {
		// calculate penetration details
		ta, tb := a.GetTransform(), b.GetTransform()
		ca, cb := a.GetCollider(), b.GetCollider()
		support := func(dir geo.Vec2) geo.Vec2 {
			return ca.Support(dir, *ta).Sub(cb.Support(dir.Negated(), *tb))
		}
		details := epa(simplex, support)
		// EPA returns the outward normal of the Minkowski difference (A-B), which points
		// from a toward b. Negate so the returned normal points from b toward a,
		// matching the convention ResolveCollision expects.
		info = CollisionInfo{A: a, B: b, Normal: details.Normal.Negated(), Depth: details.Depth}
	}
	return info, areColliding
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

func epa(simplex Simplex, support func(geo.Vec2) geo.Vec2) CollisionInfo {
	poly := NewPolytope(simplex)
	for {
		e := poly.ClosestEdge()
		p := support(e.normal)
		if p.Dot(e.normal)-e.dist < epaTolerance {
			return CollisionInfo{Normal: e.normal, Depth: e.dist}
		}
		poly.Insert(e.index, p)
	}
}
