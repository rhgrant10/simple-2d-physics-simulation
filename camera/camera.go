package camera

import (
	"github.com/rhgrant10/sim/geo"
)

type Camera struct {
	Pos    geo.Vec2
	Width  float64
	Height float64
	Scale  float64
}

func (c Camera) ToWorldPoint(s geo.Vec2) geo.Vec2 {
	return geo.Vec2{
		X: (s.X-c.Width/2)/c.Scale + c.Pos.X,
		Y: (-s.Y-c.Height/2)/c.Scale + c.Pos.Y,
	}
}

func (c Camera) ToScreenPoint(w geo.Vec2) geo.Vec2 {
	return geo.Vec2{
		X: (w.X-c.Pos.X)*c.Scale + c.Width/2,
		Y: (-w.Y-c.Pos.Y)*c.Scale + c.Height/2,
	}
}
