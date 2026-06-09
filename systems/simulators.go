package systems

import (
	"time"

	"github.com/rhgrant10/sim/world"
)

type Integrator interface {
	Step(w *world.World, dt float64)
}

type CollisionSolver interface {
	Solve(w *world.World)
}

type BasicSimulator struct {
	Integrator   Integrator
	Solver       CollisionSolver
	StepInterval time.Duration
	dtRemainder  float64
}

func (s *BasicSimulator) Step(w *world.World, dt float64) {
	s.Solver.Solve(w)
	interval := s.StepInterval.Seconds()
	dt += s.dtRemainder
	for dt > 0 {
		s.Integrator.Step(w, interval)
		dt -= interval
	}
	s.dtRemainder = interval + dt
}
