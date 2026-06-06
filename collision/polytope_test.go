package collision

import (
	"math"
	"testing"

	"github.com/rhgrant10/sim/geo"
)

func makePolytope(vertices ...geo.Vec2) Polytope {
	return Polytope{vertices: vertices}
}

func TestPolytopeNewPolytope(t *testing.T) {
	var s Simplex
	p0, p1, p2 := geo.Vec2{X: 1, Y: 0}, geo.Vec2{X: 0, Y: 1}, geo.Vec2{X: -1, Y: 0}
	s.Add(p0)
	s.Add(p1)
	s.Add(p2)

	poly := NewPolytope(s)

	if len(poly.vertices) != 3 {
		t.Fatalf("len(vertices) = %d; want 3", len(poly.vertices))
	}
	for i, want := range []geo.Vec2{p0, p1, p2} {
		if poly.vertices[i] != want {
			t.Errorf("vertices[%d] = %v; want %v", i, poly.vertices[i], want)
		}
	}
}

func TestPolytopeInsert(t *testing.T) {
	p := geo.Vec2{X: 0, Y: -1}
	base := []geo.Vec2{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}}

	tests := []struct {
		name  string
		index int
		want  []geo.Vec2
	}{
		{"front", 0, []geo.Vec2{p, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}}},
		{"middle", 2, []geo.Vec2{{X: 1, Y: 0}, {X: 0, Y: 1}, p, {X: -1, Y: 0}}},
		{"end", 3, []geo.Vec2{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}, p}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			poly := makePolytope(append([]geo.Vec2{}, base...)...)
			poly.Insert(tc.index, p)

			if len(poly.vertices) != len(tc.want) {
				t.Fatalf("len = %d; want %d", len(poly.vertices), len(tc.want))
			}
			for i, want := range tc.want {
				if poly.vertices[i] != want {
					t.Errorf("vertices[%d] = %v; want %v", i, poly.vertices[i], want)
				}
			}
		})
	}
}

// TestPolytopeClosestEdge verifies ClosestEdge using spec-derived geometric properties:
//   - dist is the perpendicular distance from the origin to the edge
//   - normal is a unit vector, perpendicular to the edge, pointing away from the origin
//   - index is the correct point to insert a new vertex between the edge's two vertices
//   - it is the minimum-distance edge in the polytope
func TestPolytopeClosestEdge(t *testing.T) {
	tests := []struct {
		name      string
		poly      Polytope
		edgeA     geo.Vec2 // first vertex of the expected closest edge
		edgeB     geo.Vec2 // second vertex of the expected closest edge
		wantIndex int
		wantDist  float64
	}{
		{
			// bottom edge (vertices[1]→vertices[2]) at y=-1 is closest; dist=1
			// other two edges have dist=4/√13 ≈ 1.109
			name:      "middle edge",
			poly:      makePolytope(geo.Vec2{X: 0, Y: 2}, geo.Vec2{X: -2, Y: -1}, geo.Vec2{X: 2, Y: -1}),
			edgeA:     geo.Vec2{X: -2, Y: -1},
			edgeB:     geo.Vec2{X: 2, Y: -1},
			wantIndex: 2,
			wantDist:  1,
		},
		{
			// bottom edge (vertices[2]→vertices[0], the wrap-around) at y=-1 is closest; dist=1
			// other edges have dist=8/√29 ≈ 1.486 and 12/√34 ≈ 2.057
			name:      "wrap-around edge",
			poly:      makePolytope(geo.Vec2{X: 2, Y: -1}, geo.Vec2{X: 0, Y: 4}, geo.Vec2{X: -3, Y: -1}),
			edgeA:     geo.Vec2{X: -3, Y: -1},
			edgeB:     geo.Vec2{X: 2, Y: -1},
			wantIndex: 0,
			wantDist:  1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := tc.poly.ClosestEdge()

			if e.index != tc.wantIndex {
				t.Errorf("index = %d; want %d", e.index, tc.wantIndex)
			}
			if math.Abs(e.dist-tc.wantDist) > eps {
				t.Errorf("dist = %v; want %v", e.dist, tc.wantDist)
			}

			// normal must be a unit vector
			if math.Abs(e.normal.Magnitude()-1) > eps {
				t.Errorf("normal magnitude = %v; want 1", e.normal.Magnitude())
			}

			// normal must be perpendicular to the edge
			edgeVec := tc.edgeB.Sub(tc.edgeA)
			if math.Abs(e.normal.Dot(edgeVec)) > eps {
				t.Errorf("normal %v not perpendicular to edge %v", e.normal, edgeVec)
			}

			// normal must point away from origin: n·(point on edge) == dist
			if math.Abs(e.normal.Dot(tc.edgeA)-tc.wantDist) > eps {
				t.Errorf("normal.Dot(edgeA) = %v; want dist=%v", e.normal.Dot(tc.edgeA), tc.wantDist)
			}

			// it must be the minimum: no other edge should be closer
			verts := tc.poly.vertices
			for i := range verts {
				j := (i + 1) % len(verts)
				ab := verts[j].Sub(verts[i])
				n := geo.Vec2{X: ab.Y, Y: -ab.X}.Normalized()
				d := n.Dot(verts[i])
				if d < 0 {
					d = -d
				}
				if d < e.dist-eps {
					t.Errorf("edge %d→%d has dist=%v, less than returned dist=%v", i, j, d, e.dist)
				}
			}
		})
	}
}
