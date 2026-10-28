package systems

import (
	"math"
	"time"

	"github.com/rhgrant10/sim/collision"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
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

type StaticIntervalSimulator struct {
	Integrator   Integrator
	Solver       CollisionSolver
	StepInterval time.Duration
	dtRemainder  float64
}

func (s *StaticIntervalSimulator) Step(w *world.World, dt float64) {
	s.Solver.Solve(w)
	interval := s.StepInterval.Seconds()
	dt += s.dtRemainder
	for dt > 0 {
		s.Integrator.Step(w, interval)
		dt -= interval
	}
	s.dtRemainder = interval + dt
}

type XPBDSimulator struct {
	Detector   collision.Detector
	SubSteps   uint
	Iterations uint
	Compliance float64 // 0 is billiard, approaching 1 gives mush
}

func (s *XPBDSimulator) Step(w *world.World, dt float64) {
	entities := []collision.Collidable{}
	for c := range world.Query4[collision.Collider, geo.Transform, physics.Material, physics.RigidBody](w) {
		entities = append(entities, &ColliderAdapter{
			collider:  c.A,
			transform: c.B,
			material:  c.C,
			rigidbody: c.D,
		})
	}

	h := dt / float64(s.SubSteps)
	iH := 1 / h
	for range s.SubSteps {
		// predict
		for i := range entities {
			rb := entities[i].GetRigidBody()
			gScale := 1.0
			if rb.IsStatic() {
				gScale = 0.0
			}
			rb.Velocity = rb.Velocity.Add(physics.Gravity.Scale(h * gScale))
			t := entities[i].GetTransform()
			rb.PrevPos = t.Pos
			t.Pos = t.Pos.Add(rb.Velocity.Scale(h))
		}

		// collide & solve constraints
		contactInfos := s.Detector.Detect(entities)
		for i := range contactInfos {
			contactInfos[i].RelativeVelocity = contactInfos[i].A.GetRigidBody().Velocity.Sub(contactInfos[i].B.GetRigidBody().Velocity)
		}
		for range s.Iterations {
			for i := range contactInfos {
				s.solveContact(&contactInfos[i])
			}
		}

		// finalize velocity (calculate then adjust)
		for i := range entities {
			rb := entities[i].GetRigidBody()
			rb.Velocity = entities[i].GetTransform().Pos.Sub(rb.PrevPos).Scale(iH)
		}
		for i := range contactInfos {
			s.solveVelocity(&contactInfos[i], h)
		}
	}
}

func (s *XPBDSimulator) solveContact(info *collision.CollisionInfo) {
	invMassA := info.A.GetRigidBody().InverseMass()
	invMassB := info.B.GetRigidBody().InverseMass()
	invMassTotal := invMassA + invMassB
	if invMassTotal <= 0 {
		return
	}

	deltaLambda := info.Depth / (invMassTotal + s.Compliance)
	push := info.Normal.Scale(deltaLambda)

	info.A.GetTransform().Pos = info.A.GetTransform().Pos.Add(push.Scale(invMassA))
	info.B.GetTransform().Pos = info.B.GetTransform().Pos.Sub(push.Scale(invMassB))

	info.Depth -= deltaLambda * invMassTotal
}

func (s *XPBDSimulator) solveVelocity(info *collision.CollisionInfo, h float64) {
	invMassA := info.A.GetRigidBody().InverseMass()
	invMassB := info.B.GetRigidBody().InverseMass()
	invMassTotal := invMassA + invMassB
	if invMassTotal <= 0 {
		return
	}

	relBefore := info.RelativeVelocity
	relAfter := info.A.GetRigidBody().Velocity.Sub(info.B.GetRigidBody().Velocity)
	approachSpeedBefore := relBefore.Dot(info.Normal)
	approachSpeedAfter := relAfter.Dot(info.Normal)

	if approachSpeedBefore >= 0 {
		// objects were not even approaching each other
		return
	}

	// restitution
	restitution := (info.A.GetMaterial().Restitution + info.B.GetMaterial().Restitution) / 2
	if math.Abs(approachSpeedBefore) < 2*h*math.Abs(physics.Gravity.Y) {
		restitution = 0.0
	}
	target := approachSpeedBefore * -restitution
	delta := target - approachSpeedAfter
	impulseMagnitude := delta / invMassTotal
	correction := info.Normal.Scale(impulseMagnitude)
	info.A.GetRigidBody().Velocity = info.A.GetRigidBody().Velocity.Add(correction.Scale(invMassA))
	info.B.GetRigidBody().Velocity = info.B.GetRigidBody().Velocity.Sub(correction.Scale(invMassB))

	relAfter = info.A.GetRigidBody().Velocity.Sub(info.B.GetRigidBody().Velocity)
	relNormal := relAfter.Dot(info.Normal)
	relTangent := relAfter.Sub(info.Normal.Scale(relNormal))

	if relTangent.Magnitude() > 1e-6 {
		friction := float64(info.A.GetMaterial().Friction+info.B.GetMaterial().Friction) / 2
		t := relTangent.Normalized()
		jt := -relTangent.Dot(t) / invMassTotal
		if maxJt := friction * math.Abs(impulseMagnitude); jt > maxJt {
			jt = maxJt
		} else if minJt := -friction * math.Abs(impulseMagnitude); jt < minJt {
			jt = minJt
		}
		info.A.GetRigidBody().Velocity = info.A.GetRigidBody().Velocity.Add(t.Scale(jt * invMassA))
		info.B.GetRigidBody().Velocity = info.B.GetRigidBody().Velocity.Sub(t.Scale(jt * invMassB))
	}
}
