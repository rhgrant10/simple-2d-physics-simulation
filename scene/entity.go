package scene

import (
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/physics"
)

type Updatable interface {
	Update(dt float64)
}

type Drawable interface {
	GetBrush(c *camera.Camera) Brush
	// Draw(screen *ebiten.Image, camera *camera.Camera)
}

type Entity interface {
	Updatable
	Drawable
	physics.Collidable
}
