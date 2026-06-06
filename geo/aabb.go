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

func (a *AABB) Left() float64 {
	return a.Min.X
}

func (a *AABB) Right() float64 {
	return a.Max.X
}

func (a *AABB) Bottom() float64 {
	return a.Min.Y
}

func (a *AABB) Top() float64 {
	return a.Max.Y
}

func (a *AABB) Contains(b AABB) bool {
	return (a.Min.X <= b.Min.X &&
		a.Min.Y <= b.Min.Y &&
		a.Max.X >= b.Max.X &&
		a.Max.Y >= b.Max.Y)
}

func (a *AABB) Overlaps(b AABB) bool {
	return (overlaps(a.Min.X, a.Max.X, b.Min.X, b.Max.X) &&
		overlaps(a.Min.Y, a.Max.Y, b.Min.Y, b.Max.Y))
}

func (a *AABB) Union(b AABB) AABB {
	return NewAABB(
		Vec2{
			X: min(a.Min.X, b.Min.X),
			Y: min(a.Min.Y, b.Min.Y),
		},
		Vec2{
			X: max(a.Max.X, b.Max.X),
			Y: max(a.Max.Y, b.Max.Y),
		},
	)
}

func overlaps(minV1, maxV1, minV2, maxV2 float64) bool {
	if minV1 <= minV2 {
		if minV2 <= maxV1 {
			return true
		}
	} else {
		if minV1 <= maxV2 {
			return true
		}
	}
	return false
}
