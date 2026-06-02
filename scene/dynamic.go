package scene

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
)

type Ball struct {
	Transform geo.Transform
	RigidBody physics.RigidBody
	Radius    float64
	Color     color.Color
	Collider  physics.CircleCollider
}

func NewBall(
	pos geo.Vec2,
	vel geo.Vec2,
	radius float64,
	density float64,
	color color.Color,
	restitution float64,
) Ball {
	area := math.Pi * radius * radius
	mass := density * area
	return Ball{
		Transform: geo.Transform{Pos: pos},
		RigidBody: physics.RigidBody{Velocity: vel, Mass: mass, Restitution: restitution},
		Radius:    radius,
		Color:     color,
		Collider:  physics.CircleCollider{Radius: radius},
	}
}

func (b *Ball) GetCollider() physics.Collider {
	return &b.Collider
}

func (b *Ball) GetTransform() *geo.Transform {
	return &b.Transform
}

func (b *Ball) GetRigidBody() *physics.RigidBody {
	return &b.RigidBody
}

func (b *Ball) Update(dt float64) {
	accel := physics.Gravity // todo: dynamics
	b.Transform.Pos = b.Transform.Pos.Add(b.RigidBody.Velocity.Scale(dt)).Add(accel.Scale(dt * dt).Scale(0.5))
	b.RigidBody.Velocity = b.RigidBody.Velocity.Add(accel.Scale(dt))
}

func (b *Ball) Draw(screen *ebiten.Image, camera *camera.Camera) {
	point := camera.ToScreenPoint(b.Transform.Pos)
	radius := b.Radius * camera.Scale
	vector.FillCircle(
		screen,
		float32(point.X),
		float32(point.Y),
		float32(radius),
		b.Color,
		camera.Antialias,
	)
}
