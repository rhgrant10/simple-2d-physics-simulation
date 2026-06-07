package collision

type Resolver interface {
	Resolve(collision CollisionInfo)
}
