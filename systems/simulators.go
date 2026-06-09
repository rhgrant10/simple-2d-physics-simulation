package systems

import (
	"github.com/rhgrant10/sim/world"
)

type Integrator interface {
	Step(w *world.World, dt float64)
}

type CollisionSolver interface {
	Solve(w *world.World)
}

type BasicSimulator struct {
	Integrator Integrator
	Solver     CollisionSolver
}

func (s *BasicSimulator) Step(w *world.World, dt float64) {
	s.Solver.Solve(w)
	s.Integrator.Step(w, dt)
}
