package systems

import (
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
	"github.com/rhgrant10/sim/world"
)

type EulerianIntegrator struct{}

func (i *EulerianIntegrator) Step(w *world.World, dt float64) {
	accel := physics.Gravity
	for v := range world.Query2[geo.Transform, physics.RigidBody](w) {
		transform := v.A
		rb := v.B
		v0 := rb.Velocity
		transform.Pos = transform.Pos.Add(rb.Velocity.Scale(dt * rb.InverseMass()))
		rb.Velocity = v0.Add(accel.Scale(dt * rb.InverseMass()))
	}
}

type SymplecticEulerianIntegrator struct{}

func (i *SymplecticEulerianIntegrator) Step(w *world.World, dt float64) {
	accel := physics.Gravity
	for v := range world.Query2[geo.Transform, physics.RigidBody](w) {
		transform := v.A
		rb := v.B
		rb.Velocity = rb.Velocity.Add(accel.Scale(dt * rb.InverseMass()))
		transform.Pos = transform.Pos.Add(rb.Velocity.Scale(dt * rb.InverseMass()))
	}
}

type VerletIntegrator struct{}

func (i *VerletIntegrator) Step(w *world.World, dt float64) {
	accel := physics.Gravity
	for v := range world.Query2[geo.Transform, physics.RigidBody](w) {
		transform := v.A
		rb := v.B
		transform.Pos = transform.Pos.Add(rb.Velocity.Scale(dt)).Add(accel.Scale(dt * dt * rb.InverseMass()).Scale(0.5))
		rb.Velocity = rb.Velocity.Add(accel.Scale(dt * rb.InverseMass()))
	}
}
