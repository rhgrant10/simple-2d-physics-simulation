package collision

import (
	"math"
	"slices"

	"github.com/rhgrant10/sim/geo"
)

type edge struct {
	normal geo.Vec2
	dist   float64
	index  int
}

type Polytope struct {
	vertices []geo.Vec2
}

func NewPolytope(s Simplex) Polytope {
	vertices := make([]geo.Vec2, 0, 16) // pre-allocation for the insertions that follow
	vertices = append(vertices, s.points[0], s.points[1], s.points[2])
	return Polytope{vertices: vertices}
}

func (p *Polytope) Insert(index int, v geo.Vec2) {
	p.vertices = slices.Insert(p.vertices, index, v)
}

func (p *Polytope) ClosestEdge() edge {
	closest := edge{dist: math.Inf(1)}

	for i := range p.vertices {
		j := (i + 1) % len(p.vertices)
		a, b := p.vertices[i], p.vertices[j]

		ab := b.Sub(a)
		n := geo.Vec2{X: ab.Y, Y: -ab.X}.Normalized() // 90° rotation
		d := n.Dot(a)
		if d < 0 {
			n, d = n.Negated(), -d
		}

		if d < closest.dist {
			closest = edge{normal: n, dist: d, index: j}
		}
	}

	return closest
}
