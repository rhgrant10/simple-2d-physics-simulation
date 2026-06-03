package scene

import (
	"image/color"

	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
)

type Box struct {
	Transform geo.Transform
	RigidBody physics.RigidBody
	Width     float64
	Height    float64
	Color     color.RGBA
	Collider  physics.BoxCollider
	Brush     BoxBrush
}

func NewBox(
	center geo.Vec2,
	width float64,
	height float64,
	color color.RGBA,
	restitution float64,
	friction physics.KineticFriction,
) Box {
	return Box{
		Transform: geo.Transform{Pos: center},
		RigidBody: physics.RigidBody{Restitution: restitution, Friction: float64(friction)},
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

func (b *Box) GetBrush(camera *camera.Camera) Brush {
	b.Brush.center = camera.ToScreenPoint(b.Transform.Pos)
	b.Brush.width = b.Width * camera.Scale
	b.Brush.height = b.Height * camera.Scale
	b.Brush.color = b.Color
	return &b.Brush
}

func (b *Box) Update(dt float64) {
}
