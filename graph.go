// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// graph.go
//
// Definitions for unweighted graphs, and methods not specific to directed
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

// Ok validates that no arcs in g lead outside the graph.
//
// Ints in an adjacency list represent half arcs.  Ok
// returns true if all int values are valid slice indexes back into g.
func (g AdjacencyList) Ok() bool {
	for _, nbs := range g {
		for _, nb := range nbs {
			if nb < 0 || nb >= len(g) {
				return false
			}
		}
	}
	return true
}

// ArcSize returns the number of arcs in g.
//
// Note that for an undirected graph, the number of edges -- the traditional
// meaning of graph size -- will be m/2.
func (g AdjacencyList) ArcSize() (m int) {
	for _, to := range g {
		m += len(to)
	}
	return
}

// Copy makes a copy of g, copying the underlying slices.
// Copy also computes the size m, the number of arcs.
func (g AdjacencyList) Copy() (c AdjacencyList, m int) {
	c = make(AdjacencyList, len(g))
	for n, to := range g {
		c[n] = append([]int{}, to...)
		m += len(to)
	}
	return
}

// CopyUndir makes an undirected copy of g -- for each arc in g, both a forward
// and reciprocal arc are added to the result.
//
// The result m here indicates the number of undirected edges, where each
// edge is represented with 2 arcs.  That will be the number of arcs in the
// input graph g, but the undirected result will have 2*m arcs.
func (g AdjacencyList) CopyUndir() (c AdjacencyList, m int) {
	c, m = g.Copy()          // start with forward links
	for gFr, to := range g { // then add reciprocals
		for _, gTo := range to {
			c[gTo] = append(c[gTo], gFr)
		}
	}
	return
}

// Undirected determines if an adjacency list is undirected.
//
// An adjacency list represents an undirected graph if every arc has a
// reciprocal.  That is, that for every arc from->to, that to->from also
// exists.
//
// If the graph is undirected, Undirected returns true, -1, -1.  If an arc
// is found witout a reciprocal, Undirected returns false along with the
// from and to nodes of an arc that represents a counterexample to the graph
// being undirected.
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

func (g LabeledAdjacencyList) Undirected() (u bool, from int, to Half) {
	for from, nbs := range g {
	nb:
		for _, to := range nbs {
			for _, r := range g[to.To] {
				if r.To == from {
					continue nb
				}
			}
			return false, from, to
		}
	}
	return true, -1, Half{}
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

// IsTreeDirected identifies trees in directed graphs.
//
// IsTreeDirected returns true if the subgraph reachable from
// root is a tree.  It does not validate that the entire graph is a tree.
func (g AdjacencyList) IsTreeDirected(root int) bool {
	var v big.Int
	var df func(int) bool
	df = func(n int) bool {
		if v.Bit(n) == 1 {
			return false
		}
		v.SetBit(&v, n, 1)
		for _, to := range g[n] {
			if !df(to) {
				return false
			}
		}
		return true
	}
	return df(root)
}

// IsTreeUndirected identifies trees in undirected graphs.
//
// IsTreeUndirected returns true if the connected component
// containing argument root is a tree.  It does not validate
// that the entire graph is a tree.
func (g AdjacencyList) IsTreeUndirected(root int) bool {
	var v big.Int
	var df func(int, int) bool
	df = func(fr, n int) bool {
		if v.Bit(n) == 1 {
			return false
		}
		v.SetBit(&v, n, 1)
		for _, to := range g[n] {
			if to != fr && !df(n, to) {
				return false
			}
		}
		return true
	}
	v.SetBit(&v, root, 1)
	for _, to := range g[root] {
		if !df(root, to) {
			return false
		}
	}
	return true
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
