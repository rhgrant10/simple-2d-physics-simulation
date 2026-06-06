package collision

import (
	"math"
	"testing"

	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
)

const eps = 1e-9

// testBody is a minimal Collidable for use in collision tests.
type testBody struct {
	transform geo.Transform
	rigidBody physics.RigidBody
}

func (tb *testBody) GetCollider() Collider            { return nil }
func (tb *testBody) GetTransform() *geo.Transform     { return &tb.transform }
func (tb *testBody) GetRigidBody() *physics.RigidBody { return &tb.rigidBody }

func makeBody(px, py, vx, vy, mass, restitution float64) *testBody {
	return &testBody{
		transform: geo.Transform{Pos: geo.Vec2{X: px, Y: py}},
		rigidBody: physics.RigidBody{Velocity: geo.Vec2{X: vx, Y: vy}, Mass: mass, Restitution: restitution},
	}
}

// TestResolveCollision_Separating: already-separating bodies must not be modified.
func TestResolveCollision_Separating(t *testing.T) {
	// a moving right (+x), b moving left (-x) — separating along n=(1,0)
	a := makeBody(1, 0, 2, 0, 1, 1)
	b := makeBody(-1, 0, -2, 0, 1, 1)
	va0, vb0 := a.rigidBody.Velocity, b.rigidBody.Velocity

	Resolve(a, b, CollisionInfo{Normal: geo.Vec2{X: 1, Y: 0}, Depth: 0.1})

	if a.rigidBody.Velocity != va0 || b.rigidBody.Velocity != vb0 {
		t.Errorf("separating bodies were modified: va=%v vb=%v; want %v %v",
			a.rigidBody.Velocity, b.rigidBody.Velocity, va0, vb0)
	}
}

// TestResolveCollision_MomentumConserved: total linear momentum must be unchanged.
// Spec: ma·va' + mb·vb' == ma·va + mb·vb (componentwise)
func TestResolveCollision_MomentumConserved(t *testing.T) {
	tests := []struct {
		name string
		a, b *testBody
	}{
		// Equal masses m=1 — the easy case
		{"equal mass m=1 elastic", makeBody(1, 0, -1, 0, 1, 1), makeBody(-1, 0, 1, 0, 1, 1)},
		{"equal mass m=1 inelastic", makeBody(1, 0, -2, 0, 1, 0), makeBody(-1, 0, 0, 0, 1, 0)},
		// Unequal masses — exposes the impulse-without-mass-scaling bug
		{"unequal mass 2:1", makeBody(1, 0, 0, 0, 2, 1), makeBody(-1, 0, 2, 0, 1, 1)},
		{"unequal mass 1:3", makeBody(1, 0, -3, 0, 1, 1), makeBody(-1, 0, 1, 0, 3, 1)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ma, mb := tc.a.rigidBody.Mass, tc.b.rigidBody.Mass
			pBefore := tc.a.rigidBody.Velocity.Scale(ma).Add(tc.b.rigidBody.Velocity.Scale(mb))
			info := CollisionInfo{Normal: tc.a.transform.Pos.Sub(tc.b.transform.Pos).Normalized(), Depth: 0.1}

			Resolve(tc.a, tc.b, info)

			pAfter := tc.a.rigidBody.Velocity.Scale(ma).Add(tc.b.rigidBody.Velocity.Scale(mb))

			if math.Abs(pAfter.X-pBefore.X) > eps || math.Abs(pAfter.Y-pBefore.Y) > eps {
				t.Errorf("momentum not conserved: before=%v after=%v", pBefore, pAfter)
			}
		})
	}
}

// TestResolveCollision_Restitution: relative velocity along n must flip by factor e.
// Spec: (va' - vb')·n == -e · (va - vb)·n
func TestResolveCollision_Restitution(t *testing.T) {
	tests := []struct {
		name string
		a, b *testBody
	}{
		{"elastic e=1", makeBody(1, 0, -1, 0, 1, 1), makeBody(-1, 0, 1, 0, 1, 1)},
		{"inelastic e=0", makeBody(1, 0, -2, 0, 1, 0), makeBody(-1, 0, 0, 0, 1, 0)},
		{"partial e=0.5", makeBody(1, 0, -1, 0, 1, 0.5), makeBody(-1, 0, 1, 0, 1, 0.5)},
		// Unequal masses
		{"unequal mass 2:1 elastic", makeBody(1, 0, 0, 0, 2, 1), makeBody(-1, 0, 2, 0, 1, 1)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := tc.a.transform.Pos.Sub(tc.b.transform.Pos).Normalized()
			e := (tc.a.rigidBody.Restitution + tc.b.rigidBody.Restitution) / 2
			relBefore := tc.a.rigidBody.Velocity.Sub(tc.b.rigidBody.Velocity).Dot(n)

			Resolve(tc.a, tc.b, CollisionInfo{Normal: n, Depth: 0.1})

			relAfter := tc.a.rigidBody.Velocity.Sub(tc.b.rigidBody.Velocity).Dot(n)
			want := -e * relBefore
			if math.Abs(relAfter-want) > eps {
				t.Errorf("restitution wrong: (va'-vb')·n = %v; want %v (e=%.1f, relBefore=%v)",
					relAfter, want, e, relBefore)
			}
		})
	}
}

// TestResolveCollision_ImpulseAlongNormal: the impulse must be purely along n.
// Spec: velocity components perpendicular to n are unchanged for both bodies.
func TestResolveCollision_ImpulseAlongNormal(t *testing.T) {
	// n = (0,1); tangential direction is x — x-velocities must not change
	a := makeBody(0, 1, -1, -3, 1, 1)
	b := makeBody(0, -1, 2, 4, 1, 1)
	vax0, vbx0 := a.rigidBody.Velocity.X, b.rigidBody.Velocity.X

	Resolve(a, b, CollisionInfo{Normal: geo.Vec2{X: 0, Y: 1}, Depth: 0.1})

	if math.Abs(a.rigidBody.Velocity.X-vax0) > eps {
		t.Errorf("a tangential velocity changed: got %v; want %v", a.rigidBody.Velocity.X, vax0)
	}
	if math.Abs(b.rigidBody.Velocity.X-vbx0) > eps {
		t.Errorf("b tangential velocity changed: got %v; want %v", b.rigidBody.Velocity.X, vbx0)
	}
}

// TestResolveCollision_StaticBody: a body with mass=0 (infinite mass) must not move.
func TestResolveCollision_StaticBody(t *testing.T) {
	// a approaches the static b; only a's velocity should change
	a := makeBody(1, 0, -2, 0, 1, 1)
	b := makeBody(-1, 0, 0, 0, 0, 1) // mass=0 → static
	vb0 := b.rigidBody.Velocity

	Resolve(a, b, CollisionInfo{Normal: geo.Vec2{X: 1, Y: 0}, Depth: 0.1})

	if b.rigidBody.Velocity != vb0 {
		t.Errorf("static body velocity changed: got %v; want %v", b.rigidBody.Velocity, vb0)
	}
}
