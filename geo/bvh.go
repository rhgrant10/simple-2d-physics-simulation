package geo

import (
	"math"
	"slices"
)

type bvnode struct {
	index    int
	aabb     AABB
	children []*bvnode
}

type BVH struct {
	root  *bvnode
	leafs map[int]*bvnode
}

func NewBVH() BVH {
	return BVH{leafs: make(map[int]*bvnode)}
}

func (b *BVH) Build(aabbs []AABB) {
	if len(aabbs) == 0 {
		panic("Cannot build a BVH with empty set of AABBs")
	}
	leafs := make([]bvnode, 0, len(aabbs))
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for i, aabb := range aabbs {
		minX = min(minX, aabb.Min.X)
		minY = min(minY, aabb.Min.Y)
		maxX = max(maxX, aabb.Max.X)
		maxY = max(maxY, aabb.Max.Y)
		node := bvnode{index: i, aabb: aabb}
		leafs = append(leafs, node)
	}
	longAxisIsX := (maxX - minX) > (maxY - minY)
	root := buildBVH(leafs, longAxisIsX)
	b.root = &root
	for i := range leafs {
		b.leafs[leafs[i].index] = &leafs[i]
	}
}

func buildBVH(nodes []bvnode, longAxisIsX bool) bvnode {
	if len(nodes) == 1 {
		return nodes[0]
	}
	if longAxisIsX {
		slices.SortFunc(nodes, func(a, b bvnode) int {
			delta := a.aabb.Center().X - b.aabb.Center().X
			if delta == 0 {
				return 0
			} else if delta < 0 {
				return -1
			}
			return 1
			// return cmp.Compare(a.aabb.Center().X, b.aabb.Center().X)
		})
	} else {
		slices.SortFunc(nodes, func(a, b bvnode) int {
			delta := a.aabb.Center().Y - b.aabb.Center().Y
			if delta == 0 {
				return 0
			} else if delta < 0 {
				return -1
			}
			return 1
			// return cmp.Compare(a.aabb.Center().Y, b.aabb.Center().Y)
		})
	}
	mid := len(nodes) / 2
	children := []*bvnode{}
	var aabb *AABB
	if len(nodes[:mid]) > 0 {
		rightChild := buildBVH(nodes[:mid], !longAxisIsX)
		aabb = &rightChild.aabb
		children = append(children, &rightChild)
	}
	if len(nodes[mid:]) > 0 {
		leftChild := buildBVH(nodes[mid:], !longAxisIsX)
		if aabb == nil {
			aabb = &leftChild.aabb
		} else {
			combinedAABB := aabb.Union(leftChild.aabb)
			aabb = &combinedAABB
		}
		children = append(children, &leftChild)
	}
	return bvnode{
		index:    -1, // internal
		aabb:     *aabb,
		children: children,
	}
}

func (b *BVH) CandidatePairs(out []PairKey) []PairKey {
	out = out[:0]
	for i := range b.leafs {
		out = append(out, b.candidatePairsFor(i)...)
	}
	slices.Sort(out)
	out = slices.Compact(out)
	return out
}

func (b *BVH) candidatePairsFor(i int) []PairKey {
	result := []PairKey{}
	if b.root != nil {
		stack := []*bvnode{b.root}
		for len(stack) > 0 {
			node := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if node.aabb.Overlaps(b.leafs[i].aabb) {
				if node.index >= 0 && node.index != i {
					result = append(result, MakePair(node.index, i))
				} else {
					stack = append(stack, node.children...)
				}
			}
		}
	}
	return result
}

func (b *BVH) Clear() {
	b.root = nil
	for k := range b.leafs {
		delete(b.leafs, k)
	}
}
