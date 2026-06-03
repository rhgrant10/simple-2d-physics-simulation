package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rhgrant10/sim/geo"
)

func addQuad(
	vertices *[]ebiten.Vertex,
	indices *[]uint16,
	x0, y0, x1, y1 float64,
	u0, v0, u1, v1 float32,
	col color.RGBA,
) {
	// uint16 index means max 65535 vertices per draw call.
	// At 4 verts/shape that's ~16000 shapes — fine for 1000 balls.
	// If you ever need more, split into multiple DrawTriangles calls.
	base := uint16(len(*vertices))

	cr := float32(col.R) / 255
	cg := float32(col.G) / 255
	cb := float32(col.B) / 255
	ca := float32(col.A) / 255

	*vertices = append(*vertices,
		ebiten.Vertex{DstX: float32(x0), DstY: float32(y0), SrcX: u0, SrcY: v0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		ebiten.Vertex{DstX: float32(x1), DstY: float32(y0), SrcX: u1, SrcY: v0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		ebiten.Vertex{DstX: float32(x0), DstY: float32(y1), SrcX: u0, SrcY: v1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		ebiten.Vertex{DstX: float32(x1), DstY: float32(y1), SrcX: u1, SrcY: v1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
	)
	*indices = append(*indices,
		base+0, base+1, base+2,
		base+1, base+3, base+2,
	)
}

type Brush interface {
	AppendVertices(vertices *[]ebiten.Vertex, indices *[]uint16)
}

type CircleBrush struct {
	center geo.Vec2
	radius float64
	color  color.RGBA
}

func (b *CircleBrush) AppendVertices(vertices *[]ebiten.Vertex, indices *[]uint16) {
	addQuad(
		vertices, indices,
		b.center.X-b.radius, b.center.Y-b.radius, // top-left
		b.center.X+b.radius, b.center.Y+b.radius, // bottom-right
		0, 0, 1, 1, // UV: 0..1 signals "circle" to shader
		b.color,
	)
}

type BoxBrush struct {
	center        geo.Vec2
	width, height float64
	color         color.RGBA
}

func (b *BoxBrush) AppendVertices(vertices *[]ebiten.Vertex, indices *[]uint16) {
	addQuad(
		vertices, indices,
		b.center.X-b.width/2, b.center.Y-b.height/2,
		b.center.X+b.width/2, b.center.Y+b.height/2,
		2, 2, 3, 3, // UV: 2..3 signals "box" to shader
		b.color,
	)
}
