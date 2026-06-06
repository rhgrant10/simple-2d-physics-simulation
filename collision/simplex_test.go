package collision

import (
	"testing"

	"github.com/rhgrant10/sim/geo"
)

func TestSimplexAdd(t *testing.T) {
	var s Simplex
	p1, p2 := geo.Vec2{X: 1, Y: 2}, geo.Vec2{X: 3, Y: 4}

	s.Add(p1)
	if s.size != 1 || s.points[0] != p1 {
		t.Errorf("after first Add: size=%d points[0]=%v; want size=1 points[0]=%v", s.size, s.points[0], p1)
	}

	s.Add(p2)
	if s.size != 2 || s.points[1] != p2 {
		t.Errorf("after second Add: size=%d points[1]=%v; want size=2 points[1]=%v", s.size, s.points[1], p2)
	}
}

// Line: origin is beside the segment (ab·ao > 0).
// GJK requires dir to be perpendicular to ab and point toward the origin.
// Geometry: b=(-1,1) a=(1,1) — horizontal line above origin.
func TestSimplexContainsLine_PerpendicularSide(t *testing.T) {
	b, a := geo.Vec2{X: -1, Y: 1}, geo.Vec2{X: 1, Y: 1}
	ab, ao := b.Sub(a), a.Negated()

	var s Simplex
	s.Add(b)
	s.Add(a)

	dir := geo.Vec2{X: 0, Y: 1}
	got := s.Contains(&dir)

	if got {
		t.Error("Contains = true; want false for a line simplex")
	}
	if s.size != 2 {
		t.Errorf("size = %d; want 2", s.size)
	}
	if ab.Dot(dir) != 0 {
		t.Errorf("dir %v is not perpendicular to ab %v", dir, ab)
	}
	if ao.Dot(dir) <= 0 {
		t.Errorf("dir %v does not point toward origin (ao=%v)", dir, ao)
	}
}

// Line: a is the closest vertex to the origin (ab·ao ≤ 0).
// GJK requires the simplex to shrink to {a} and dir to equal ao.
// Geometry: b=(-3,0) a=(-2,0) — both on the negative x-axis.
func TestSimplexContainsLine_AClosest(t *testing.T) {
	b, a := geo.Vec2{X: -3, Y: 0}, geo.Vec2{X: -2, Y: 0}
	ao := a.Negated()

	var s Simplex
	s.Add(b)
	s.Add(a)

	dir := geo.Vec2{X: 0, Y: 1}
	got := s.Contains(&dir)

	if got {
		t.Error("Contains = true; want false for a line simplex")
	}
	if s.size != 1 {
		t.Errorf("size = %d; want 1", s.size)
	}
	if s.points[0] != a {
		t.Errorf("points[0] = %v; want a=%v", s.points[0], a)
	}
	if dir != ao {
		t.Errorf("dir = %v; want ao=%v", dir, ao)
	}
}

// Triangle: origin is inside.
// GJK requires Contains to return true.
// Triangle c=(2,-1) b=(-2,-1) a=(0,2) wraps the origin.
func TestSimplexContainsTriangle_ContainsOrigin(t *testing.T) {
	var s Simplex
	s.Add(geo.Vec2{X: 2, Y: -1})  // c
	s.Add(geo.Vec2{X: -2, Y: -1}) // b
	s.Add(geo.Vec2{X: 0, Y: 2})   // a

	dir := geo.Vec2{X: 0, Y: 1}
	got := s.Contains(&dir)

	if !got {
		t.Error("Contains = false; want true (origin is inside the triangle)")
	}
}

// Triangle: origin is outside the AB edge.
// GJK requires: simplex shrinks to {b,a}, dir is perpendicular to ab,
// points away from c, and points toward the origin.
// Geometry: c=(-2,6) b=(2,6) a=(0,3) — all vertices well above the origin.
func TestSimplexContainsTriangle_OutsideAB(t *testing.T) {
	c, b, a := geo.Vec2{X: -2, Y: 6}, geo.Vec2{X: 2, Y: 6}, geo.Vec2{X: 0, Y: 3}
	ab, ac, ao := b.Sub(a), c.Sub(a), a.Negated()

	var s Simplex
	s.Add(c)
	s.Add(b)
	s.Add(a)

	dir := geo.Vec2{X: 0, Y: 1}
	got := s.Contains(&dir)

	if got {
		t.Error("Contains = true; want false (origin is outside AB)")
	}
	if s.size != 2 {
		t.Errorf("size = %d; want 2", s.size)
	}
	if s.points[0] != b || s.points[1] != a {
		t.Errorf("simplex = {%v,%v}; want {b=%v, a=%v}", s.points[0], s.points[1], b, a)
	}
	if ab.Dot(dir) != 0 {
		t.Errorf("dir %v is not perpendicular to ab %v", dir, ab)
	}
	if ac.Dot(dir) >= 0 {
		t.Errorf("dir %v does not point away from c (ac·dir=%v; want <0)", dir, ac.Dot(dir))
	}
	if ao.Dot(dir) <= 0 {
		t.Errorf("dir %v does not point toward origin (ao·dir=%v; want >0)", dir, ao.Dot(dir))
	}
}

// Triangle: origin is outside the AC edge.
// GJK requires: simplex shrinks to {c,a}, dir is perpendicular to ac,
// points away from b, and points toward the origin.
// Geometry: c=(1,2) b=(3,0) a=(1,0) — AC is the vertical line x=1; origin is left of it.
func TestSimplexContainsTriangle_OutsideAC(t *testing.T) {
	c, b, a := geo.Vec2{X: 1, Y: 2}, geo.Vec2{X: 3, Y: 0}, geo.Vec2{X: 1, Y: 0}
	ab, ac, ao := b.Sub(a), c.Sub(a), a.Negated()

	var s Simplex
	s.Add(c)
	s.Add(b)
	s.Add(a)

	dir := geo.Vec2{X: 0, Y: 1}
	got := s.Contains(&dir)

	if got {
		t.Error("Contains = true; want false (origin is outside AC)")
	}
	if s.size != 2 {
		t.Errorf("size = %d; want 2", s.size)
	}
	if s.points[0] != c || s.points[1] != a {
		t.Errorf("simplex = {%v,%v}; want {c=%v, a=%v}", s.points[0], s.points[1], c, a)
	}
	if ac.Dot(dir) != 0 {
		t.Errorf("dir %v is not perpendicular to ac %v", dir, ac)
	}
	if ab.Dot(dir) >= 0 {
		t.Errorf("dir %v does not point away from b (ab·dir=%v; want <0)", dir, ab.Dot(dir))
	}
	if ao.Dot(dir) <= 0 {
		t.Errorf("dir %v does not point toward origin (ao·dir=%v; want >0)", dir, ao.Dot(dir))
	}
}
