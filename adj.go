// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"math/big"
)

// file ed.go contains definitions common to different search functions

// An AdjacencyList represents a graph as a list of neighbors for each node.
// The "node ID" of a node is simply it's slice index in the AdjacencyList.
//
// Adjacency lists are inherently directed. To represent an undirected graph,
// create reciprocal neighbors.
type AdjacencyList [][]int

func (g AdjacencyList) Undirected() (u bool, from, to int) {
	for from, nbs := range g {
	nb:
		for _, to := range nbs {
			for _, r := range g[to] {
				if r == from {
					continue nb
				}
			}
			return false, from, to
		}
	}
	return true, -1, -1
}

func (g AdjacencyList) ConnectedComponents() []int {
	var r []int
	var c big.Int
	var df func(int)
	df = func(n int) {
		c.SetBit(&c, n, 1)
		for _, nb := range g[n] {
			if c.Bit(nb) == 0 {
				df(nb)
			}
		}
	}
	for n := range g {
		if c.Bit(n) == 0 {
			r = append(r, n)
			df(n)
		}
	}
	return r
}

func (g AdjacencyList) Bipartite(n int) (b bool, c1, c2 *big.Int, oc []int) {
	c1 = &big.Int{}
	c2 = &big.Int{}
	b = true
	var open bool
	var df func(n int, c1, c2 *big.Int)
	df = func(n int, c1, c2 *big.Int) {
		c1.SetBit(c1, n, 1)
		for _, nb := range g[n] {
			if c1.Bit(nb) == 1 {
				b = false
				oc = []int{nb, n}
				open = true
				return
			}
			if c2.Bit(nb) == 1 {
				continue
			}
			df(nb, c2, c1)
			if b {
				continue
			}
			switch {
			case !open:
			case n == oc[0]:
				open = false
			default:
				oc = append(oc, n)
			}
			return
		}
	}
	df(n, c1, c2)
	if b {
		return b, c1, c2, nil
	}
	return b, nil, nil, oc
}

// ValidTo validates that no arcs in the reciever graph lead outside the graph.
//
// Ints in the adjacency list structure represent "to" half arcs.  ValidTo
// returns true if all int values are valid slice indexes back into g.
func (g AdjacencyList) ValidTo() bool {
	for _, nbs := range g {
		for _, nb := range nbs {
			if nb < 0 || nb >= len(g) {
				return false
			}
		}
	}
	return true
}

// A FromTree represents a spanning tree where each node is associated with
// a half arc identifying an arc from another node.
//
// Other terms for this data structure include "predecessor list", "in-tree",
// "inverse arborescence", and "spaghetti stack."  It is an effecient
// representation for accumulating path results for various algorithms that
// search or traverse graphs starting from a single source or start node.
//
// For a node n, Paths[n] contains information about the path from the
// the start node to n.  For reached nodes, the PathEnd.Len field will be
// > 0 and indicate the length of the path from start.  The PathEnd.From
// field will indicate the node this node was reached from, or -1 in the
// case of the start node.  For unreached nodes, PathEnd.Len will be 0 and
// PathEnd.From will be -1.
type FromTree struct {
	Start  int       // start node, argument to search, root of the tree
	Paths  []PathEnd // tree representation
	MaxLen int       // length of longest path, max of all PathEnd.Len values
}

// NewFromTree creates a FromTree object.  You don't typically call this
// function from application code.  Rather it is typically called by search
// object constructors.  NewFromTree leaves the result object with zero values
// and does not call the Reset method.
func NewFromTree(n int) *FromTree {
	return &FromTree{Paths: make([]PathEnd, n)}
}

// Reset initializes a FromTree in preparation for a search.  Search methods
// will call this function and you don't typically call it from application
// code.
func (t *FromTree) Reset() {
	t.Start = -1
	for n := range t.Paths {
		t.Paths[n] = PathEnd{From: -1, Len: 0}
	}
	t.MaxLen = 0
}

// A PathEnd associates a half arc and a path length.
//
// See FromTree for use by search functions.
type PathEnd struct {
	From int // a "from" half arc, the node the arc comes from
	Len  int // number of nodes in path from start
}

// PathTo decodes a FromTree, recovering a found path.
//
// PathTo returns the path recorded by some search from the start node of the
// search to the given end node.  The returned slice represents a sequence of
// half arcs.  If the search did not find a path  end the slice result will be
// nil.
func (t *FromTree) PathTo(end int) []int {
	n := t.Paths[end].Len
	if n == 0 {
		return nil
	}
	p := make([]int, n)
	for {
		n--
		p[n] = end
		if n == 0 {
			return p
		}
		end = t.Paths[end].From
	}
}
