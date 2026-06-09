package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/mesh"
	"github.com/rhgrant10/sim/world"
)

type BasicRenderer struct {
	vertices   []ebiten.Vertex
	indices    []uint16
	whitePixel *ebiten.Image
	shader     *ebiten.Shader
}

func NewBasicRenderer() *BasicRenderer {
	r := &BasicRenderer{}

	// 1×1 white image used as a dummy texture source
	r.whitePixel = ebiten.NewImage(1, 1)
	r.whitePixel.Fill(color.White)

	// Compile the Kage shader. This happens once at startup.
	var err error
	r.shader, err = ebiten.NewShader(shaderSrc)
	if err != nil {
		panic(err)
	}

	// TODO: we really need to preallocate?
	r.vertices = make([]ebiten.Vertex, 0, 5000)
	r.indices = make([]uint16, 0, 7500)

	return r
}

func (r *BasicRenderer) Draw(screen *ebiten.Image, camera *camera.Camera, w *world.World) {
	r.begin()
	for c := range world.Query2[geo.Transform, mesh.Circle](w) {
		r.drawCircle(camera, c.A, c.B)
	}
	for c := range world.Query2[geo.Transform, mesh.Rectangle](w) {
		r.drawRectangle(camera, c.A, c.B)
	}
	r.flush(screen)
}

func (r *BasicRenderer) drawCircle(camera *camera.Camera, t *geo.Transform, m *mesh.Circle) {
	center := camera.ToScreenPoint(t.Pos)
	radius := m.Radius * camera.Scale
	r.addQuad(
		center.X-radius, center.Y-radius, // top-left
		center.X+radius, center.Y+radius, // bottom-right
		0, 0, 1, 1, // UV: 0..1 signals "circle" to shader
		m.Color,
	)
}

func (r *BasicRenderer) drawRectangle(camera *camera.Camera, t *geo.Transform, m *mesh.Rectangle) {
	center := camera.ToScreenPoint(t.Pos)
	width := m.Width * camera.Scale
	height := m.Height * camera.Scale
	r.addQuad(
		center.X-width/2, center.Y-height/2, // top-left
		center.X+width/2, center.Y+height/2, // bottom-right
		2, 2, 3, 3, // UV: 2..3 signals "rect" to shader
		m.Color,
	)
}

func (r *BasicRenderer) begin() {
	r.vertices = r.vertices[:0]
	r.indices = r.indices[:0]
}

func (r *BasicRenderer) flush(screen *ebiten.Image) {
	if len(r.vertices) == 0 {
		return
	}

	op := &ebiten.DrawTrianglesShaderOptions{}
	// No uniforms needed: color comes from vertex data,
	// shape type comes from UV range.

	screen.DrawTrianglesShader(r.vertices, r.indices, r.shader, op)
}

func (r *BasicRenderer) addQuad(
	x0, y0, x1, y1 float64,
	u0, v0, u1, v1 float32,
	col color.RGBA,
) {
	// uint16 index means max 65535 vertices per draw call.
	// At 4 verts/shape that's ~16000 shapes — fine for 1000 balls.
	// If you ever need more, split into multiple DrawTriangles calls.
	base := uint16(len(r.vertices))

	cr := float32(col.R) / 255
	cg := float32(col.G) / 255
	cb := float32(col.B) / 255
	ca := float32(col.A) / 255

	r.vertices = append(r.vertices,
		ebiten.Vertex{DstX: float32(x0), DstY: float32(y0), SrcX: u0, SrcY: v0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		ebiten.Vertex{DstX: float32(x1), DstY: float32(y0), SrcX: u1, SrcY: v0, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		ebiten.Vertex{DstX: float32(x0), DstY: float32(y1), SrcX: u0, SrcY: v1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
		ebiten.Vertex{DstX: float32(x1), DstY: float32(y1), SrcX: u1, SrcY: v1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca},
	)
	r.indices = append(r.indices,
		base+0, base+1, base+2,
		base+1, base+3, base+2,
	)
}
