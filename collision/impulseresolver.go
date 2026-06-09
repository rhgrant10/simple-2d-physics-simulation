package collision

import (
	"math"
)

type ImpulseResolver struct {
	percent    float64
	slop       float64
	iterations int
}

func NewImpulseResolver() ImpulseResolver {
	return ImpulseResolver{percent: 0.4, slop: 0.01, iterations: 1}
}

func (r ImpulseResolver) Resolve(info CollisionInfo) {
	a := info.A
	b := info.B
	n := info.Normal
	inverseMassSum := a.GetRigidBody().InverseMass() + b.GetRigidBody().InverseMass()

	av, bv := a.GetRigidBody().Velocity, b.GetRigidBody().Velocity
	relVelVec := av.Sub(bv)
	relVel := relVelVec.Dot(n)

	if relVel < 0 {
		restitution := (a.GetMaterial().Restitution + b.GetMaterial().Restitution) / 2
		if inverseMassSum > 0 {
			// apply the collision impulse
			j := -(1 + restitution) * relVel / inverseMassSum
			impulse := n.Scale(j)
			a.GetRigidBody().Velocity = av.Add(impulse.Scale(a.GetRigidBody().InverseMass()))
			b.GetRigidBody().Velocity = bv.Sub(impulse.Scale(b.GetRigidBody().InverseMass()))

			// handle friction
			noSpinFactor := 2.0 / 7.0
			t := relVelVec.Sub(n.Scale(relVel))
			if t.Magnitude() > 1e-6 {
				t = t.Normalized()
				relVelT := relVelVec.Dot(t)
				mu := noSpinFactor * float64(a.GetMaterial().Friction+b.GetMaterial().Friction) / 2
				jt := -relVelT / inverseMassSum
				// Coulomb cone
				if jt > mu*j {
					jt = mu * j
				} else if jt < -mu*j {
					jt = -mu * j
				}
				frictionImpulse := t.Scale(jt)
				a.GetRigidBody().Velocity = a.GetRigidBody().Velocity.Add(frictionImpulse.Scale(a.GetRigidBody().InverseMass()))
				b.GetRigidBody().Velocity = b.GetRigidBody().Velocity.Sub(frictionImpulse.Scale(b.GetRigidBody().InverseMass()))
			}
		}
	}

	// Positional correction: push overlapping bodies apart to prevent sinking.
	if inverseMassSum > 0 {
		for range r.iterations {
			corrMag := math.Max(info.Depth-r.slop, 0) / inverseMassSum * r.percent
			corr := n.Scale(corrMag)
			a.GetTransform().Pos = a.GetTransform().Pos.Add(corr.Scale(a.GetRigidBody().InverseMass()))
			b.GetTransform().Pos = b.GetTransform().Pos.Sub(corr.Scale(b.GetRigidBody().InverseMass()))
		}
	}
}
