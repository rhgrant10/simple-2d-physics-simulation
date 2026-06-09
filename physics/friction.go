package physics

type KineticFriction float64

const (
	FrictionIce      KineticFriction = 0.03
	FrictionSteel    KineticFriction = 0.15
	FrictionWood     KineticFriction = 0.3
	FrictionConcrete KineticFriction = 0.6
	FrictionRubber   KineticFriction = 0.7
)
