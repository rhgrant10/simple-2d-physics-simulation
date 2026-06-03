package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	Vertices   []ebiten.Vertex
	Indices    []uint16
	whitePixel *ebiten.Image
	shader     *ebiten.Shader
}

func NewRenderer() *Renderer {
	r := &Renderer{}

	// 1×1 white image used as a dummy texture source
	r.whitePixel = ebiten.NewImage(1, 1)
	r.whitePixel.Fill(color.White)

	// Compile the Kage shader. This happens once at startup.
	var err error
	r.shader, err = ebiten.NewShader(shaderSrc)
	if err != nil {
		panic(err)
	}

	// Pre-allocate slices. 1000 balls × 4 verts each = 4000 verts.
	// Over-allocating is fine; we reset length each frame.
	r.Vertices = make([]ebiten.Vertex, 0, 5000)
	r.Indices = make([]uint16, 0, 7500)

	return r
}

func (r *Renderer) Begin() {
	r.Vertices = r.Vertices[:0]
	r.Indices = r.Indices[:0]
}

func (r *Renderer) Flush(screen *ebiten.Image) {
	if len(r.Vertices) == 0 {
		return
	}

	op := &ebiten.DrawTrianglesShaderOptions{}
	// No uniforms needed: color comes from vertex data,
	// shape type comes from UV range.

	screen.DrawTrianglesShader(r.Vertices, r.Indices, r.shader, op)
}

// addQuad is the internal helper that does the actual vertex/index work.
// It appends 4 vertices and 6 indices (two triangles) per shape.
//
//	(x0,y0)-----(x1,y0)
//	   |    \       |
//	   |      \     |
//	(x0,y1)-----(x1,y1)
//
// Triangle 1: top-left, top-right, bottom-left
// Triangle 2: top-right, bottom-right, bottom-left
func (r *Renderer) addQuad(
	x0, y0, x1, y1 float64,
	u0, v0, u1, v1 float32,
	col color.RGBA,
) {
	// uint16 index means max 65535 vertices per draw call.
	// At 4 verts/shape that's ~16000 shapes — fine for 1000 balls.
	// If you ever need more, split into multiple DrawTriangles calls.
	base := uint16(len(r.Vertices))

	cr := float32(col.R) / 255
	cg := float32(col.G) / 255
	cb := float32(col.B) / 255
	ca := float32(col.A) / 255

	r.Vertices = append(r.Vertices,
		ebiten.Vertex{DstX: float32(x0), DstY: float32(y0), SrcX: u0, SrcY: v0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		ebiten.Vertex{DstX: float32(x1), DstY: float32(y0), SrcX: u1, SrcY: v0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		ebiten.Vertex{DstX: float32(x0), DstY: float32(y1), SrcX: u0, SrcY: v1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		ebiten.Vertex{DstX: float32(x1), DstY: float32(y1), SrcX: u1, SrcY: v1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
	)
	r.Indices = append(r.Indices,
		base+0, base+1, base+2,
		base+1, base+3, base+2,
	)
}
