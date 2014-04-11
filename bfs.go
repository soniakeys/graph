package ed

import (
	"math/big"
)

type BreadthFirst struct {
	Graph  AdjacencyList
	Result *FromTree
}

func NewBreadthFirst(g AdjacencyList) *BreadthFirst {
	return &BreadthFirst{
		Graph:  g,
		Result: NewFromTree(len(g)),
	}
}

func (b *BreadthFirst) Path(start, end int) []int {
	b.Traverse(start, func(n int) bool { return n != end })
	return b.Result.PathTo(end)
}

func (b *BreadthFirst) AllPaths(start int) int {
	return b.Traverse(start, func(int) bool { return true })
}

type Visitor func(n int) (ok bool)

func (b *BreadthFirst) Traverse(start int, v Visitor) int {
	b.Result.Reset()
	rp := b.Result.Paths
	b.Result.Start = start
	level := 1
	rp[start].Len = level
	if !v(start) {
		b.Result.MaxLen = level
		return -1
	}
	nReached := 1 // accumulated for a return value
	// the frontier consists of nodes all at the same level
	frontier := []int{start}
	for {
		level++
		var next []int
		for _, n := range frontier {
			for _, nb := range b.Graph[n] {
				if rp[nb].Len == 0 {
					rp[nb] = PathEnd{From: n, Len: level}
					if !v(nb) {
						b.Result.MaxLen = level
						return -1
					}
					next = append(next, nb)
					nReached++
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
	}
	b.Result.MaxLen = level - 1
	return nReached
}

type BreadthFirst2 struct {
	To, From AdjacencyList
	M        int
	Result   *FromTree
}

func NewBreadthFirst2(to, from AdjacencyList, m int) *BreadthFirst2 {
	return &BreadthFirst2{
		To:     to,
		From:   from,
		M:      m,
		Result: NewFromTree(len(to)),
	}
}

func (g AdjacencyList) Inverse() (from AdjacencyList, m int) {
	from = make([][]int, len(g))
	for n, nbs := range g {
		for _, nb := range nbs {
			from[nb] = append(from[nb], n)
			m++
		}
	}
	return
}

func (b *BreadthFirst2) Path(start, end int) []int {
	b.Traverse(start, func(n int) bool { return n != end })
	return b.Result.PathTo(end)
}

func (b *BreadthFirst2) AllPaths(start int) int {
	return b.Traverse(start, func(int) bool { return true })
}

func (b *BreadthFirst2) Traverse(start int, v Visitor) int {
	b.Result.Reset()
	rp := b.Result.Paths
	b.Result.Start = start
	level := 1
	rp[start].Len = level
	if !v(start) {
		b.Result.MaxLen = level
		return -1
	}
	nReached := 1 // accumulated for a return value
	// the frontier consists of nodes all at the same level
	frontier := []int{start}
	mf := len(b.To[start])      // number of arcs leading out from frontier
	ctb := b.M / 10             // threshold change from top-down to bottom-up
	k14 := 14 * b.M / len(b.To) // 14 * mean degree
	cbt := len(b.To) / k14      // threshold change from bottom-up to top-down
	var fBits, nextb big.Int
	for {
		// top down step
		level++
		var next []int
		for _, n := range frontier {
			for _, nb := range b.To[n] {
				if rp[nb].Len == 0 {
					rp[nb] = PathEnd{From: n, Len: level}
					if !v(nb) {
						b.Result.MaxLen = level
						return -1
					}
					next = append(next, nb)
					nReached++
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
		if mf > ctb {
			// switch to bottom up!
		} else {
			// stick with top down
			continue
		}
		// convert frontier representation
		nf := 0 // number of vertices on the frontier
		for _, n := range frontier {
			fBits.SetBit(&fBits, n, 1)
			nf++
		}
	bottomUpLoop:
		level++
		nNext := 0
		for n := range b.From {
			if rp[n].Len == 0 {
				for _, nb := range b.From[n] {
					if fBits.Bit(nb) == 1 {
						rp[n] = PathEnd{From: nb, Len: level}
						if !v(nb) {
							b.Result.MaxLen = level
							return -1
						}
						nextb.SetBit(&nextb, n, 1)
						nReached++
						nNext++
						break
					}
				}
			}
		}
		if nNext == 0 {
			break
		}
		fBits, nextb = nextb, fBits
		nextb.SetInt64(0)
		nf = nNext
		if nf < cbt {
			// switch back to top down!
		} else {
			// stick with bottom up
			goto bottomUpLoop
		}
		// convert frontier representation
		mf = 0
		frontier = frontier[:0]
		for n := range b.To {
			if fBits.Bit(n) == 1 {
				frontier = append(frontier, n)
				mf += len(b.To[n])
			}
		}
		fBits.SetInt64(0)
	}
	b.Result.MaxLen = level - 1
	return nReached
}
