package geo

type AABB struct {
	Min Vec2
	Max Vec2
}

func NewAABB(p1, p2 Vec2) AABB {
	if p1.X < p2.X {
		return AABB{Min: p1, Max: p2}
	}
	return AABB{Min: p2, Max: p1}
}
