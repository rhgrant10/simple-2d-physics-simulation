package collision

import "github.com/rhgrant10/sim/geo"

type Broadphaser interface {
	Build(bvs *[]geo.AABB)
	CandidatePairs(out []geo.PairKey) []geo.PairKey
	Clear()
}

type NoBroadphase struct {
	indices map[int]struct{}
}

func NewNoBroadphase() NoBroadphase {
	return NoBroadphase{indices: make(map[int]struct{})}
}

func (b *NoBroadphase) Build(bvs *[]geo.AABB) {
	for i := range *bvs {
		b.indices[i] = struct{}{}
	}
}

func (b *NoBroadphase) CandidatePairs(out []geo.PairKey) []geo.PairKey {
	out = out[:0]
	for a := range b.indices {
		for b := range b.indices {
			if a != b {
				out = append(out, geo.MakePair(a, b))
			}
		}
	}
	return out
}

func (b *NoBroadphase) Clear() {
	for k := range b.indices {
		delete(b.indices, k)
	}
}
