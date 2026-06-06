package geo

import "log"

type bvnode struct {
	index    int
	bv       AABB
	children []*bvnode
}

type BVH struct {
	root             *bvnode
	insertionCounter int
}

func (b *BVH) Insert(i int, aabb AABB) {
	b.insertionCounter++
	n := &bvnode{index: i, bv: aabb, children: []*bvnode{}}
	if b.root == nil {
		b.root = n
		return
	}

	// find the first place where n.bv is within parent.bv but
	// also doesn't overlap with any bv's of its children
	stack := []*bvnode{b.root}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		withinAChild := false
		for c := range curr.children {
			if curr.children[c].bv.Contains(n.bv) {
				// is in the parent and also within this child. must be a child
				withinAChild = true
				stack = append(stack, curr.children[c])
				break
			} else if curr.children[c].bv.Overlaps(n.bv) {
				// is in the parent and partially in this child. must be a sibling
				curr.children = append(curr.children, n)
				return
			}
		}
		if !withinAChild {
			// its in the parent, but not in any children. must be a sibling
			curr.children = append(curr.children, n)
			return
		}
	}

	// isn't within the root!
	if n.bv.Contains(b.root.bv) {
		// n contains root
		n.children = append(n.children, b.root)
		b.root = n
	} else {
		// overlaps or disjoint - new parent to contain them both
		b.root = &bvnode{
			bv:       n.bv.Union(b.root.bv),
			children: []*bvnode{n, b.root},
		}
	}
}

func (b *BVH) CandidatePairs(out []PairKey) []PairKey {
	out = out[:0]
	visitCount := 0
	if b.root != nil {
		stack := []*bvnode{b.root}
		for len(stack) > 0 {
			visitCount++
			curr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			log.Printf("Node has %d children", len(curr.children))
			if len(curr.children) == 2 && curr.children[0].bv.Overlaps(curr.children[1].bv) {
				out = append(out, MakePair(curr.children[0].index, curr.children[1].index))
			}
			stack = append(stack, curr.children...)
		}
	}
	log.Printf("Visited %d nodes while finding CandidatePairs", visitCount)
	return out
}

func (b *BVH) Clear() {
	log.Printf("Clearing BVH (%d nodes)", b.insertionCounter)
	b.root = nil
	b.insertionCounter = 0
}
