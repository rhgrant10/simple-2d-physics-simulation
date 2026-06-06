package geo

import "testing"

func TestAABBContains(t *testing.T) {
	testCases := []struct {
		desc     string
		a        AABB
		b        AABB
		expected bool
	}{
		{
			desc: "Strict contains",
			a: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: -5, Y: -5},
				Max: Vec2{X: 5, Y: 5},
			},
			expected: true,
		},
		{
			desc: "Equal contains",
			a: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: 10, Y: 10},
			},
			expected: true,
		},
		{
			desc: "Disjoint",
			a: AABB{
				Min: Vec2{X: 1, Y: 1},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: -1, Y: -1},
			},
			expected: false,
		},
		{
			desc: "Overlap",
			a: AABB{
				Min: Vec2{X: -1, Y: -1},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: 1, Y: 1},
			},
			expected: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.a.Contains(tC.b)
			if actual != tC.expected {
				t.Errorf("got %v, want %v", actual, tC.expected)
			}
		})
	}
}

func TestAABBOverlaps(t *testing.T) {
	testCases := []struct {
		desc     string
		a        AABB
		b        AABB
		expected bool
	}{
		{
			desc: "Strict overlaps",
			a: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: -5, Y: -5},
				Max: Vec2{X: 5, Y: 5},
			},
			expected: true,
		},
		{
			desc: "Equal overlaps",
			a: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: 10, Y: 10},
			},
			expected: true,
		},
		{
			desc: "Disjoint",
			a: AABB{
				Min: Vec2{X: 1, Y: 1},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: -1, Y: -1},
			},
			expected: false,
		},
		{
			desc: "Overlap",
			a: AABB{
				Min: Vec2{X: -1, Y: -1},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: -10, Y: -10},
				Max: Vec2{X: 1, Y: 1},
			},
			expected: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.a.Overlaps(tC.b)
			if actual != tC.expected {
				t.Errorf("got %v, want %v", actual, tC.expected)
			}
		})
	}
}

func TestAABBUnion(t *testing.T) {
	testCases := []struct {
		desc     string
		a        AABB
		b        AABB
		expected AABB
	}{
		{
			desc: "Containing",
			a: AABB{
				Min: Vec2{X: 0, Y: 0},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: 1, Y: 1},
				Max: Vec2{X: 9, Y: 9},
			},
			expected: AABB{
				Min: Vec2{X: 0, Y: 0},
				Max: Vec2{X: 10, Y: 10},
			},
		},
		{
			desc: "Overlapping",
			a: AABB{
				Min: Vec2{X: 0, Y: 0},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: 5, Y: 5},
				Max: Vec2{X: 15, Y: 15},
			},
			expected: AABB{
				Min: Vec2{X: 0, Y: 0},
				Max: Vec2{X: 15, Y: 15},
			},
		},
		{
			desc: "Disjoint",
			a: AABB{
				Min: Vec2{X: 0, Y: 0},
				Max: Vec2{X: 10, Y: 10},
			},
			b: AABB{
				Min: Vec2{X: 20, Y: 20},
				Max: Vec2{X: 30, Y: 30},
			},
			expected: AABB{
				Min: Vec2{X: 0, Y: 0},
				Max: Vec2{X: 30, Y: 30},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.a.Union(tC.b)
			if actual != tC.expected {
				t.Errorf("got %v, want %v", actual, tC.expected)
			}
		})
	}
}
