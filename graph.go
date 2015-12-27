// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// graph.go
//
// Definitions for unlabeled graphs, and methods not specific to directed
// or undirected graphs.  Method docs need not mention that they work on both.

package graph

import (
	"math/big"
	"sort"
)

var one = big.NewInt(1)

// An AdjacencyList represents a graph as a list of neighbors for each node.
// The "node ID" of a node is simply it's slice index in the AdjacencyList.
//
// Adjacency lists are inherently directed. To represent an undirected graph,
// create reciprocal neighbors.
type AdjacencyList [][]int

// BoundsOk validates that all arcs in g stay within the slice bounds of g.
//
// BoundsOk returns true when no arcs point outside the bounds of g.
// Otherwise it returns false and an example arc that points outside of g.
func (g AdjacencyList) BoundsOk() (ok bool, fr, to int) {
	for fr, to := range g {
		for _, to := range to {
			if to < 0 || to >= len(g) {
				return false, fr, to
			}
		}
	}
	return true, -1, -1
}

// ArcSize returns the number of arcs in g.
//
// Note that for an undirected graph witout loops, the number of edges --
// the traditional meaning of graph size -- will be m/2.
func (g AdjacencyList) ArcSize() (m int) {
	for _, to := range g {
		m += len(to)
	}
	return
}

// Copy makes a copy of g, copying the underlying slices.
// Copy also computes the arc size m, the number of arcs.
func (g AdjacencyList) Copy() (c AdjacencyList, m int) {
	c = make(AdjacencyList, len(g))
	for n, to := range g {
		c[n] = append([]int{}, to...)
		m += len(to)
	}
	return
}

// Simple checks for loops and parallel arcs.
//
// A graph is "simple" if it has no loops or parallel arcs.
//
// Simple returns true, -1 for simple graphs.  If a loop or parallel arc is
// found, simple returns false and and a node that represents a counterexample
// to the graph being simple.
func (g AdjacencyList) Simple() (s bool, n int) {
	var t []int
	for n, nbs := range g {
		if len(nbs) == 0 {
			continue
		}
		t = append(t[:0], nbs...)
		sort.Ints(t)
		if t[0] == n {
			return false, n
		}
		for i, nb := range t[1:] {
			if nb == n || nb == t[i] {
				return false, n
			}
		}
	}
	return true, -1
}

// DepthFirst traverses a graph depth first.
//
// As it traverses it calls visitor function v for each node.  If v returns
// false at any point, the traversal is terminated immediately and DepthFirst
// returns false.  Otherwise DepthFirst returns true.
//
// DepthFirst uses argument bm is used as a bitmap to guide the traversal.
// For a complete traversal, bm should be 0 initially.  During the
// traversal, bits are set corresponding to each node visited.
// The bit is set before calling the visitor function.
//
// Argument bm can be nil if you have no need for it.
// In this case a bitmap is created internally for one-time use.
//
// Alternatively v can be nil.  In this case traversal still procedes and
// updates the bitmap, which can be a useful result.
// DepthFirst always returns true in this case.
//
// It makes no sense for both bm and v to be nil.  In this case DepthFirst
// returns false immediately.
func (g AdjacencyList) DepthFirst(start int, bm *big.Int, v Visitor) (ok bool) {
	if bm == nil {
		if v == nil {
			return false
		}
		bm = new(big.Int)
	}
	ok = true
	var df func(n int)
	df = func(n int) {
		if bm.Bit(n) == 1 {
			return
		}
		bm.SetBit(bm, n, 1)
		if v != nil && !v(n) {
			ok = false
			return
		}
		for _, nb := range g[n] {
			df(nb)
		}
	}
	df(start)
	return
}
