package geo

import (
	"math/rand"
	"testing"
)

// bruteForceOverlaps returns every pair (i,j) where i<j and aabbs[i] overlaps aabbs[j].
func bruteForceOverlaps(aabbs []AABB) map[PairKey]struct{} {
	pairs := make(map[PairKey]struct{})
	for i := 0; i < len(aabbs); i++ {
		for j := i + 1; j < len(aabbs); j++ {
			if aabbs[i].Overlaps(aabbs[j]) {
				pairs[MakePair(i, j)] = struct{}{}
			}
		}
	}
	return pairs
}

func buildAndQuery(aabbs []AABB) []PairKey {
	bvh := NewBVH()
	bvh.Build(aabbs)
	return bvh.CandidatePairs(nil)
}

// centered returns an AABB centered at (cx,cy) with the given half-extents.
func centered(cx, cy, hw, hh float64) AABB {
	return AABB{
		Min: Vec2{X: cx - hw, Y: cy - hh},
		Max: Vec2{X: cx + hw, Y: cy + hh},
	}
}

// TestBVHCompleteness verifies that the BVH never misses a pair of overlapping AABBs.
// False positives (extra candidate pairs) are acceptable; false negatives are bugs.
func TestBVHCompleteness(t *testing.T) {
	tests := []struct {
		name  string
		aabbs []AABB
	}{
		{
			name:  "single",
			aabbs: []AABB{centered(0, 0, 1, 1)},
		},
		{
			name:  "two non-overlapping",
			aabbs: []AABB{centered(0, 0, 1, 1), centered(10, 0, 1, 1)},
		},
		{
			name:  "two overlapping",
			aabbs: []AABB{centered(0, 0, 1, 1), centered(1, 0, 1, 1)},
		},
		{
			name: "three in a line, only adjacent overlap",
			// bodies 0-1 overlap, 1-2 overlap, 0-2 do not
			aabbs: []AABB{
				centered(0, 0, 1.5, 1),
				centered(2, 0, 1.5, 1),
				centered(4, 0, 1.5, 1),
			},
		},
		{
			name: "input order reversed relative to sort order",
			// bodies are stored right-to-left, so after the BVH sorts by X,
			// every body ends up at a different position than its index
			aabbs: []AABB{
				centered(10, 0, 1, 1), // body 0 – rightmost
				centered(7, 0, 1, 1),  // body 1
				centered(4, 0, 1, 1),  // body 2
				centered(1, 0, 1, 1),  // body 3 – leftmost
			},
		},
		{
			name: "two clusters, input order reversed",
			// bodies 0 and 1 overlap each other in the right cluster;
			// bodies 2 and 3 are in the left cluster and don't overlap anyone.
			// Input order is right-to-left so sorting swaps indices vs positions.
			aabbs: []AABB{
				centered(11, 0, 1, 1), // body 0, right cluster
				centered(10, 0, 1, 1), // body 1, right cluster – overlaps body 0
				centered(1, 0, 1, 1),  // body 2, left cluster
				centered(6, 0, 1, 1),  // body 3, middle – no overlaps
			},
		},
		{
			name: "all identical AABBs",
			aabbs: []AABB{
				centered(0, 0, 1, 1),
				centered(0, 0, 1, 1),
				centered(0, 0, 1, 1),
				centered(0, 0, 1, 1),
			},
		},
		{
			name: "5x5 grid, no overlaps",
			aabbs: func() []AABB {
				var out []AABB
				for r := 0; r < 5; r++ {
					for c := 0; c < 5; c++ {
						out = append(out, centered(float64(c)*4, float64(r)*4, 1, 1))
					}
				}
				return out
			}(),
		},
		{
			name: "row of 10 touching circles",
			aabbs: func() []AABB {
				var out []AABB
				for i := 0; i < 10; i++ {
					out = append(out, centered(float64(i)*2, 0, 1, 1))
				}
				return out
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ground := bruteForceOverlaps(tc.aabbs)
			candidates := make(map[PairKey]struct{})
			for _, p := range buildAndQuery(tc.aabbs) {
				candidates[p] = struct{}{}
			}
			for pair := range ground {
				if _, ok := candidates[pair]; !ok {
					t.Errorf("missed overlapping pair (%d, %d)", pair.A(), pair.B())
				}
			}
		})
	}
}

// TestBVHCompletenessRandom runs randomised trials comparing the BVH against
// the brute-force reference. A fixed seed keeps the test deterministic.
func TestBVHCompletenessRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for trial := 0; trial < 50; trial++ {
		n := 2 + rng.Intn(60)
		aabbs := make([]AABB, n)
		for i := range aabbs {
			cx := rng.Float64() * 200
			cy := rng.Float64() * 200
			r := 1 + rng.Float64()*8
			aabbs[i] = centered(cx, cy, r, r)
		}

		ground := bruteForceOverlaps(aabbs)
		candidates := make(map[PairKey]struct{})
		for _, p := range buildAndQuery(aabbs) {
			candidates[p] = struct{}{}
		}

		for pair := range ground {
			if _, ok := candidates[pair]; !ok {
				t.Errorf("trial %d (n=%d): missed overlapping pair (%d, %d)",
					trial, n, pair.A(), pair.B())
			}
		}
	}
}

func TestBVHEmptyDoesPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Build must panic on empty input, but didn't")
		}
	}()
	buildAndQuery(nil)
}

// TestBVHNoSelfPairs verifies that no body is ever paired with itself.
func TestBVHNoSelfPairs(t *testing.T) {
	aabbs := []AABB{
		centered(0, 0, 1, 1),
		centered(0, 0, 1, 1), // identical to 0 – guaranteed overlap with everyone
		centered(5, 0, 1, 1),
	}
	for _, p := range buildAndQuery(aabbs) {
		if p.A() == p.B() {
			t.Errorf("self-pair returned: (%d, %d)", p.A(), p.B())
		}
	}
}

// TestBVHNoDuplicatePairs verifies that each candidate pair appears at most once.
func TestBVHNoDuplicatePairs(t *testing.T) {
	aabbs := []AABB{
		centered(0, 0, 1, 1),
		centered(0, 0, 1, 1),
		centered(0, 0, 1, 1),
	}
	seen := make(map[PairKey]struct{})
	for _, p := range buildAndQuery(aabbs) {
		if _, dup := seen[p]; dup {
			t.Errorf("duplicate pair: (%d, %d)", p.A(), p.B())
		}
		seen[p] = struct{}{}
	}
}
