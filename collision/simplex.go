package collision

import "github.com/rhgrant10/sim/geo"

type Simplex struct {
	points [3]geo.Vec2
	size   int
}

func (s *Simplex) Add(p geo.Vec2) {
	s.points[s.size] = p
	s.size++
}

func (s *Simplex) Contains(dir *geo.Vec2) bool {
	switch s.size {
	case 2:
		return s.handleLine(dir)
	case 3:
		return s.handleTriangle(dir)
	}
	return false
}

func (s *Simplex) handleLine(dir *geo.Vec2) bool {
	b, a := s.points[0], s.points[1]
	ab, ao := b.Sub(a), a.Negated()

	if ab.Dot(ao) > 0 {
		*dir = ab.TripleProduct(ao, ab)
	} else {
		s.points[0] = a
		s.size = 1
		*dir = ao
	}
	return false
}

func (s *Simplex) handleTriangle(dir *geo.Vec2) bool {
	c, b, a := s.points[0], s.points[1], s.points[2]
	ab, ac, ao := b.Sub(a), c.Sub(a), a.Negated()
	abPerp := ac.TripleProduct(ab, ab)
	acPerp := ab.TripleProduct(ac, ac)

	if abPerp.Dot(ao) > 0 {
		s.points = [3]geo.Vec2{b, a}
		s.size = 2
		*dir = abPerp
		return false
	}
	if acPerp.Dot(ao) > 0 {
		s.points = [3]geo.Vec2{c, a}
		s.size = 2
		*dir = acPerp
		return false
	}
	return true
}
