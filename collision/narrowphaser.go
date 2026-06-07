package collision

type NarrowPhaser interface {
	Check(a, b Collidable) (info CollisionInfo, areColliding bool)
}
