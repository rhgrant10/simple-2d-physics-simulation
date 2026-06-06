package scene

import (
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/collision"
)

type Updatable interface {
	Update(dt float64)
}

type Drawable interface {
	GetBrush(c *camera.Camera) Brush
}

type Entity interface {
	Updatable
	Drawable
	collision.Collidable
}
