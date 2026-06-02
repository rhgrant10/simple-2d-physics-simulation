package physics

import (
	"math"
	"testing"

	"github.com/rhgrant10/sim/geo"
)

// dcEps is a looser tolerance for DetectCollision tests: EPA approximates
// curved shapes (circles) with polygon edges, so normals won't be exact.
const dcEps = 1e-3

type circleBody struct {
	transform geo.Transform
	collider  CircleCollider
	rigidBody RigidBody
}

func (c *circleBody) GetCollider() Collider        { return &c.collider }
func (c *circleBody) GetTransform() *geo.Transform { return &c.transform }
func (c *circleBody) GetRigidBody() *RigidBody     { return &c.rigidBody }

func makeCircleBody(px, py, radius float64) *circleBody {
	return &circleBody{
		transform: geo.Transform{Pos: geo.Vec2{X: px, Y: py}},
		collider:  CircleCollider{Radius: radius},
	}
}

// TestDetectCollision_NoOverlap: separated circles must return false.
func TestDetectCollision_NoOverlap(t *testing.T) {
	a := makeCircleBody(5, 0, 1)
	b := makeCircleBody(-5, 0, 1)

	_, ok := DetectCollision(a, b)

	if ok {
		t.Error("expected no collision, got one")
	}
}

// TestDetectCollision_YAxis: circles overlapping along y.
// Depth = r1+r2 - dist = 2 - 1.4 = 0.6; Normal points from b toward a = (0,+1).
//
// NOTE: axis-aligned x-penetration (posA.Y == posB.Y) is degenerate with the
// current GJK starting direction of (1,0) and should not be used as a test case.
func TestDetectCollision_YAxis(t *testing.T) {
	a := makeCircleBody(0, 0.7, 1)
	b := makeCircleBody(0, -0.7, 1)

	info, ok := DetectCollision(a, b)

	if !ok {
		t.Fatal("expected collision, got none")
	}
	if math.Abs(info.Depth-0.6) > dcEps {
		t.Errorf("Depth = %v; want 0.6", info.Depth)
	}
	if math.Abs(info.Normal.X) > dcEps || math.Abs(info.Normal.Y-1) > dcEps {
		t.Errorf("Normal = %v; want (0, 1)", info.Normal)
	}
}

// TestDetectCollision_Symmetry: swapping a and b must negate the Normal; Depth unchanged.
func TestDetectCollision_Symmetry(t *testing.T) {
	a := makeCircleBody(0, 0.7, 1)
	b := makeCircleBody(0, -0.7, 1)

	ab, okAB := DetectCollision(a, b)
	ba, okBA := DetectCollision(b, a)

	if !okAB || !okBA {
		t.Fatal("expected both to detect collision")
	}
	if math.Abs(ab.Depth-ba.Depth) > dcEps {
		t.Errorf("depths differ: %v vs %v", ab.Depth, ba.Depth)
	}
	if math.Abs(ab.Normal.X+ba.Normal.X) > dcEps || math.Abs(ab.Normal.Y+ba.Normal.Y) > dcEps {
		t.Errorf("normals not negated: ab=%v ba=%v", ab.Normal, ba.Normal)
	}
}
