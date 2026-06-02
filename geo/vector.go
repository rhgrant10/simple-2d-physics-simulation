package geo

import "math"

type Vec2 struct {
	X float64
	Y float64
}

func Vec2FromPolar(r, deg float64) Vec2 {
	return Vec2{
		X: math.Cos(math.Pi*deg/180) * r,
		Y: math.Sin(math.Pi*deg/180) * r,
	}
}

func (a Vec2) Add(b Vec2) Vec2 {
	return Vec2{X: a.X + b.X, Y: a.Y + b.Y}
}

func (a Vec2) Sub(b Vec2) Vec2 {
	return Vec2{X: a.X - b.X, Y: a.Y - b.Y}
}

func (a Vec2) Scale(b float64) Vec2 {
	return Vec2{X: a.X * b, Y: a.Y * b}
}

func (a Vec2) Magnitude() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

func (a Vec2) IsZero() bool {
	return a.X == 0 && a.Y == 0
}

func (a Vec2) Negated() Vec2 {
	return Vec2{X: -a.X, Y: -a.Y}
}

func (a Vec2) Dot(b Vec2) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vec2) Normalized() Vec2 {
	m := a.Magnitude()
	return Vec2{X: a.X / m, Y: a.Y / m}
}

func (a Vec2) TripleProduct(b, c Vec2) Vec2 {
	// 2D version: (a × b) is a scalar, then rotate c by 90° scaled by it
	ab := a.X*b.Y - a.Y*b.X
	return Vec2{-ab * c.Y, ab * c.X}
}
