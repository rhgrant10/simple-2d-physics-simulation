package geo

import "math"

type cellKey struct{ x, y int }
type PairKey struct{ A, B int }

func makePair(i, j int) PairKey {
	if i < j {
		return PairKey{i, j}
	}
	return PairKey{j, i}
}

type SpatialHash struct {
	cellSize    float64
	invCellSize float64
	cells       map[cellKey][]int
	pairSet     map[PairKey]struct{}
}

func NewSpatialHash(cellSize float64, cellCount, pairCount int) *SpatialHash {
	return &SpatialHash{
		cellSize:    cellSize,
		invCellSize: 1 / cellSize,
		cells:       make(map[cellKey][]int, cellCount),
		pairSet:     make(map[PairKey]struct{}, pairCount),
	}
}

func (sh *SpatialHash) componentToCell(v float64) int {
	return int(math.Floor(v * sh.invCellSize))
}

func (sh *SpatialHash) Insert(i int, aabb AABB) {
	x0 := sh.componentToCell(aabb.Min.X)
	x1 := sh.componentToCell(aabb.Min.Y)
	y0 := sh.componentToCell(aabb.Max.X)
	y1 := sh.componentToCell(aabb.Max.Y)

	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			k := cellKey{x, y}
			sh.cells[k] = append(sh.cells[k], i)
		}
	}
}

// CandidatePairs returns all pairs (i,j) that share at least one cell.
// The returned slice is reused across calls – copy it if you need to keep it.
func (sh *SpatialHash) CandidatePairs(out []PairKey) []PairKey {
	out = out[:0]
	for _, occupants := range sh.cells {
		// All occupants of the same cell are candidate pairs.
		for a := 0; a < len(occupants)-1; a++ {
			for b := a + 1; b < len(occupants); b++ {
				pk := makePair(occupants[a], occupants[b])
				if _, seen := sh.pairSet[pk]; !seen {
					sh.pairSet[pk] = struct{}{}
					out = append(out, pk)
				}
			}
		}
	}
	return out
}

// Clear resets the hash for the next frame.
// It reuses existing map buckets to avoid GC pressure.
func (sh *SpatialHash) Clear() {
	for k, s := range sh.cells {
		sh.cells[k] = s[:0] // keep backing array, zero length
	}
	for k := range sh.pairSet {
		delete(sh.pairSet, k)
	}
}
