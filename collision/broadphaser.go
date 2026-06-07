package collision

import "github.com/rhgrant10/sim/geo"

type Broadphaser interface {
	Build(bvs *[]geo.AABB)
	CandidatePairs(out []geo.PairKey) []geo.PairKey
	Clear()
}
