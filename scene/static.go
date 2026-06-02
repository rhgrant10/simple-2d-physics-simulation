package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
)

type Box struct {
	Transform geo.Transform
	RigidBody physics.RigidBody
	Width     float64
	Height    float64
	Color     color.Color
	Collider  physics.BoxCollider
}

func NewBox(
	center geo.Vec2,
	width float64,
	height float64,
	color color.Color,
	restitution float64,
) Box {
	return Box{
		Transform: geo.Transform{Pos: center},
		RigidBody: physics.RigidBody{Restitution: restitution},
		Width:     width,
		Height:    height,
		Color:     color,
		Collider:  physics.BoxCollider{HalfW: width / 2, HalfH: height / 2},
	}
}

func (b *Box) GetCollider() physics.Collider {
	return &b.Collider
}

func (b *Box) GetTransform() *geo.Transform {
	return &b.Transform
}

func (b *Box) GetRigidBody() *physics.RigidBody {
	return &b.RigidBody
}

func (b *Box) Update(dt float64) {
	// b.RigidBody.Velocity.X = 0
	// b.RigidBody.Velocity.Y = 0
}

func (b *Box) Draw(screen *ebiten.Image, camera *camera.Camera) {
	point := camera.ToScreenPoint(
		geo.Vec2{
			X: b.Transform.Pos.X - b.Width/2,
			Y: b.Transform.Pos.Y + b.Height/2,
		},
	)
	width := b.Width * camera.Scale
	height := b.Height * camera.Scale
	vector.FillRect(
		screen,
		float32(point.X),
		float32(point.Y),
		float32(width),
		float32(height),
		b.Color,
		camera.Antialias,
	)
}
