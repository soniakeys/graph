// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import "math/big"

// same as BreadthFirst
type BreadthFirstBit struct {
	Graph  AdjacencyList
	Result *FromTree
}

func NewBreadthFirstBit(g AdjacencyList) *BreadthFirstBit {
	return &BreadthFirstBit{
		Graph:  g,
		Result: newFromTree(len(g)),
	}
}

func (b *BreadthFirstBit) AllPaths(start int) int {
	return b.Traverse(start, func(int) bool { return true })
}

func (b *BreadthFirstBit) Traverse(start int, v Visitor) int {
	b.Result.reset()
	var vis big.Int
	rp := b.Result.Paths
	level := 1
	vis.SetBit(&vis, start, 1)
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
				//				if rp[nb].Len == 0 {
				if vis.Bit(nb) == 0 {
					vis.SetBit(&vis, nb, 1)
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
