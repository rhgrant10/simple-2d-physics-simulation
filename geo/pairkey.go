package geo

type PairKey uint64

func MakePair(a, b int) PairKey {
	if a > b {
		a, b = b, a
	}
	return PairKey(uint64(uint32(a))<<32 | uint64(uint32(b)))
}

func (pk PairKey) A() int {
	return int(pk >> 32)
}

func (pk PairKey) B() int {
	return int(uint32(pk))
}
