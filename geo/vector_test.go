package geo

import (
	"math"
	"testing"
)

func TestVec2Add(t *testing.T) {
	tests := []struct {
		a, b, want Vec2
	}{
		{Vec2{1, 2}, Vec2{3, 4}, Vec2{4, 6}},
		{Vec2{0, 0}, Vec2{5, -3}, Vec2{5, -3}},
		{Vec2{-1, -2}, Vec2{1, 2}, Vec2{0, 0}},
	}
	for _, tc := range tests {
		got := tc.a.Add(tc.b)
		if got != tc.want {
			t.Errorf("Add(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestVec2Sub(t *testing.T) {
	tests := []struct {
		a, b, want Vec2
	}{
		{Vec2{5, 7}, Vec2{2, 3}, Vec2{3, 4}},
		{Vec2{0, 0}, Vec2{1, 1}, Vec2{-1, -1}},
		{Vec2{3, 3}, Vec2{3, 3}, Vec2{0, 0}},
	}
	for _, tc := range tests {
		got := tc.a.Sub(tc.b)
		if got != tc.want {
			t.Errorf("Sub(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestVec2Scale(t *testing.T) {
	tests := []struct {
		a    Vec2
		s    float64
		want Vec2
	}{
		{Vec2{2, 3}, 2, Vec2{4, 6}},
		{Vec2{1, -1}, -1, Vec2{-1, 1}},
		{Vec2{5, 5}, 0, Vec2{0, 0}},
	}
	for _, tc := range tests {
		got := tc.a.Scale(tc.s)
		if got != tc.want {
			t.Errorf("Scale(%v, %v) = %v, want %v", tc.a, tc.s, got, tc.want)
		}
	}
}

func TestVec2Magnitude(t *testing.T) {
	tests := []struct {
		a    Vec2
		want float64
	}{
		{Vec2{3, 4}, 5},
		{Vec2{0, 0}, 0},
		{Vec2{1, 0}, 1},
		{Vec2{0, -7}, 7},
	}
	for _, tc := range tests {
		got := tc.a.Magnitude()
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("Magnitude(%v) = %v, want %v", tc.a, got, tc.want)
		}
	}
}

func TestVec2IsZero(t *testing.T) {
	tests := []struct {
		a    Vec2
		want bool
	}{
		{Vec2{0, 0}, true},
		{Vec2{1, 0}, false},
		{Vec2{0, 1}, false},
		{Vec2{-1, -1}, false},
	}
	for _, tc := range tests {
		got := tc.a.IsZero()
		if got != tc.want {
			t.Errorf("IsZero(%v) = %v, want %v", tc.a, got, tc.want)
		}
	}
}

func TestVec2Negated(t *testing.T) {
	tests := []struct {
		a, want Vec2
	}{
		{Vec2{1, 2}, Vec2{-1, -2}},
		{Vec2{0, 0}, Vec2{0, 0}},
		{Vec2{-3, 4}, Vec2{3, -4}},
	}
	for _, tc := range tests {
		got := tc.a.Negated()
		if got != tc.want {
			t.Errorf("Negated(%v) = %v, want %v", tc.a, got, tc.want)
		}
	}
}

func TestVec2Dot(t *testing.T) {
	tests := []struct {
		a, b Vec2
		want float64
	}{
		{Vec2{1, 0}, Vec2{0, 1}, 0},
		{Vec2{2, 3}, Vec2{4, 5}, 23},
		{Vec2{1, 1}, Vec2{-1, 1}, 0},
		{Vec2{3, 4}, Vec2{3, 4}, 25},
	}
	for _, tc := range tests {
		got := tc.a.Dot(tc.b)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("Dot(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestVec2Normalized(t *testing.T) {
	tests := []struct {
		a Vec2
	}{
		{Vec2{3, 4}},
		{Vec2{1, 0}},
		{Vec2{0, 5}},
		{Vec2{-2, 2}},
	}
	for _, tc := range tests {
		got := tc.a.Normalized()
		mag := got.Magnitude()
		if math.Abs(mag-1) > 1e-9 {
			t.Errorf("Normalized(%v).Magnitude() = %v, want 1", tc.a, mag)
		}
	}
}

func TestVec2TripleProduct(t *testing.T) {
	// (a × b) × c in 2D: result should be perpendicular to a
	tests := []struct {
		a, b, c, want Vec2
	}{
		{Vec2{1, 0}, Vec2{0, 1}, Vec2{1, 0}, Vec2{0, 1}},
		{Vec2{1, 0}, Vec2{1, 0}, Vec2{1, 1}, Vec2{0, 0}},
	}
	for _, tc := range tests {
		got := tc.a.TripleProduct(tc.b, tc.c)
		if math.Abs(got.X-tc.want.X) > 1e-9 || math.Abs(got.Y-tc.want.Y) > 1e-9 {
			t.Errorf("TripleProduct(%v, %v, %v) = %v, want %v", tc.a, tc.b, tc.c, got, tc.want)
		}
	}
}
