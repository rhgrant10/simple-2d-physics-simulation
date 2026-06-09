package collision

type Solver interface {
	Solve(entities []Collidable)
}

type BasicSolver struct {
	detector   Detector
	resovler   Resolver
	iterations int
}

func NewSolver(d Detector, r Resolver) BasicSolver {
	return BasicSolver{detector: d, resovler: r, iterations: 1}
}

func (h *BasicSolver) Solve(entities []Collidable) {
	collisions := h.detector.Detect(entities)
	for i := range collisions {
		h.resovler.Resolve(collisions[i])
	}
}
