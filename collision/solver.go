package collision

type Solver struct {
	detector   Detector
	resovler   Resolver
	iterations int
}

func NewSolver(d Detector, r Resolver) Solver {
	return Solver{detector: d, resovler: r, iterations: 1}
}

func (h *Solver) Solve(entities []Collidable) {
	collisions := h.detector.Detect(entities)
	for i := range collisions {
		h.resovler.Resolve(collisions[i])
	}
}
