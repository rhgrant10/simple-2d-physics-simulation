package collision

type NarrowPhaser interface {
	Detect(a, b Collidable) (info CollisionInfo, areColliding bool)
}
