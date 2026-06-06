package collision

import (
	"math"
	"testing"

	"github.com/rhgrant10/sim/geo"
)

const epaEps = 1e-6

// rectSupport returns a support function for an axis-aligned rectangle
// (Minkowski difference) centered at (cx, cy) with half-extents (hw, hh).
func rectSupport(cx, cy, hw, hh float64) func(geo.Vec2) geo.Vec2 {
	return func(d geo.Vec2) geo.Vec2 {
		x := cx + hw
		if d.X < 0 {
			x = cx - hw
		}
		y := cy + hh
		if d.Y < 0 {
			y = cy - hh
		}
		return geo.Vec2{X: x, Y: y}
	}
}

// TestEPA_XAxisPenetration: the Minkowski difference is a rectangle whose
// closest edge to the origin is the left face (x = −0.5, dist = 0.5).
// The initial simplex already touches that boundary, so EPA terminates
// after one edge query. Expected: Normal=(−1,0), Depth=0.5.
func TestEPA_XAxisPenetration(t *testing.T) {
	// Rectangle at center=(0.5,0), half-extents=(1,2).
	// Edges: left x=−0.5 (dist=0.5), right x=1.5, top/bottom y=±2.
	support := rectSupport(0.5, 0, 1, 2)

	// Three corners of the rectangle that enclose the origin.
	// Verified: origin lies strictly inside this triangle (signed-area test).
	var s Simplex
	s.Add(geo.Vec2{X: 1.5, Y: 0})   // rightmost boundary point
	s.Add(geo.Vec2{X: -0.5, Y: -2}) // bottom-left corner
	s.Add(geo.Vec2{X: -0.5, Y: 2})  // top-left corner

	info := epa(s, support)

	if math.Abs(info.Depth-0.5) > epaEps {
		t.Errorf("Depth = %v; want 0.5", info.Depth)
	}
	if math.Abs(info.Normal.X-(-1)) > epaEps || math.Abs(info.Normal.Y) > epaEps {
		t.Errorf("Normal = %v; want (-1, 0)", info.Normal)
	}
}

// TestEPA_YAxisPenetration: the Minkowski difference is a rectangle whose
// closest edge is the bottom face (y = −0.6, dist = 0.6). The initial
// simplex's closest edge is the hypotenuse (dist ≈ 0.358), so EPA must
// insert one new vertex before it finds the true boundary.
// Expected: Normal=(0,−1), Depth=0.6.
func TestEPA_YAxisPenetration(t *testing.T) {
	// Rectangle at center=(0,0.4), half-extents=(2,1).
	// Edges: bottom y=−0.6 (dist=0.6), top y=1.4, left/right x=±2.
	support := rectSupport(0, 0.4, 2, 1)

	// Three corners of the rectangle that enclose the origin.
	// Closest simplex edge is the hypotenuse (dist≈0.358), not the bottom — EPA expands.
	var s Simplex
	s.Add(geo.Vec2{X: -2, Y: -0.6}) // bottom-left corner
	s.Add(geo.Vec2{X: 2, Y: -0.6})  // bottom-right corner
	s.Add(geo.Vec2{X: 2, Y: 1.4})   // top-right corner

	info := epa(s, support)

	if math.Abs(info.Depth-0.6) > epaEps {
		t.Errorf("Depth = %v; want 0.6", info.Depth)
	}
	if math.Abs(info.Normal.X) > epaEps || math.Abs(info.Normal.Y-(-1)) > epaEps {
		t.Errorf("Normal = %v; want (0, -1)", info.Normal)
	}
}
