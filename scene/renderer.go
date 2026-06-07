package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rhgrant10/sim/camera"
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

func (r *Renderer) Render(screen *ebiten.Image, camera *camera.Camera, entities []Entity) {
	r.begin()
	for i := range entities {
		entities[i].GetBrush(camera).AppendVertices(&r.Vertices, &r.Indices)
	}
	r.flush(screen)
}

func (r *Renderer) begin() {
	r.Vertices = r.Vertices[:0]
	r.Indices = r.Indices[:0]
}

func (r *Renderer) flush(screen *ebiten.Image) {
	if len(r.Vertices) == 0 {
		return
	}

	op := &ebiten.DrawTrianglesShaderOptions{}
	// No uniforms needed: color comes from vertex data,
	// shape type comes from UV range.

	screen.DrawTrianglesShader(r.Vertices, r.Indices, r.shader, op)
}
