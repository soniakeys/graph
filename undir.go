// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// undir.go
//
// Methods specific to undirected graphs.
// Doc for each method should specifically say undirected.

package graph

import (
	"math/big"
)

// ConnectedComponents, for undirected graphs, determines the connected
// components in g.
//
// Returned is a slice with a single representative node from each connected
// component and a parallel slice with the number of nodes in the corresponding
// component.
func (g AdjacencyList) ConnectedComponents() (rep, nNodes []int) {
	var c big.Int
	var nc int
	var df func(int)
	df = func(n int) {
		c.SetBit(&c, n, 1)
		nc++
		for _, nb := range g[n] {
			if c.Bit(nb) == 0 {
				df(nb)
			}
		}
		return
	}
	for n := range g {
		if c.Bit(n) == 0 {
			rep = append(rep, n)
			nc = 0
			df(n)
			nNodes = append(nNodes, nc)
		}
	}
	return
}

// wip
type BiconnectedComponents struct {
	Graph  AdjacencyList
	Start  int
	Cuts   big.Int // bitmap of node cuts
	From   []int   // from-tree
	Leaves []int   // leaves of from-tree
}

func NewBiconnectedComponents(g AdjacencyList) *BiconnectedComponents {
	return &BiconnectedComponents{
		Graph: g,
		From:  make([]int, len(g)),
	}
}

func (b *BiconnectedComponents) Find(start int) {
	g := b.Graph
	depth := make([]int, len(g))
	low := make([]int, len(g))
	// reset from any previous run
	b.Cuts.SetInt64(0)
	bf := b.From
	for n := range bf {
		bf[n] = -1
	}
	b.Leaves = b.Leaves[:0]
	d := 1 // depth. d > 0 means visited
	depth[start] = d
	low[start] = d
	d++
	var df func(int, int)
	df = func(from, n int) {
		bf[n] = from
		depth[n] = d
		dn := d
		l := d
		d++
		cut := false
		leaf := true
		for _, nb := range g[n] {
			if depth[nb] == 0 {
				leaf = false
				df(n, nb)
				if low[nb] < l {
					l = low[nb]
				}
				if low[nb] >= dn {
					cut = true
				}
			} else if nb != from && depth[nb] < l {
				l = depth[nb]
			}
		}
		low[n] = l
		if cut {
			b.Cuts.SetBit(&b.Cuts, n, 1)
		}
		if leaf {
			b.Leaves = append(b.Leaves, n)
		}
		d--
	}
	nbs := g[start]
	if len(nbs) == 0 {
		return
	}
	df(start, nbs[0])
	var rc uint
	for _, nb := range nbs[1:] {
		if depth[nb] == 0 {
			rc = 1
			df(start, nb)
		}
	}
	b.Cuts.SetBit(&b.Cuts, start, rc)
	return
}
