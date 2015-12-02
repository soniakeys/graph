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

// FromList represents a tree where each node is associated with
// a half arc identifying an arc from another node.
//
// Paths represents a tree with information about the path to each node
// from a start node.  See PathEnd documentation.
//
// Leaves and MaxLen are used by some search and traversal functions to
// return extra information.  Where Leaves is used it serves as a bitmap where
// Leave.Bit() == 1 for each leaf of the tree.  Where MaxLen is used it is
// provided primarily as a convenience for functions that might want to
// anticipate the maximum path length that would be encountered traversing
// the tree.
//
// A single FromList can also represent a forest.  In this case paths from
// all leaves do not return to a single start node, but multiple start nodes.
type FromList struct {
	Paths  []PathEnd // tree representation
	Leaves big.Int   // leaves of tree
	MaxLen int       // length of longest path, max of all PathEnd.Len values
}

// PathEnd associates a half arc and a path length.
//
// PathEnd is an element type of FromList, a return type from various search
// functions.
//
// For a start node of a search, From will be -1 and Len will be 1. For other
// nodes reached by the search, From represents a half arc in a path back to
// start and Len represents the number of nodes in the path.  For nodes not
// reached by the search, From will be -1 and Len will be 0.
type PathEnd struct {
	From int // a "from" half arc, the node the arc comes from
	Len  int // number of nodes in path from start
}

// NewFromList creates a FromList object.  You don't typically call this
// function from application code.  Rather it is typically called by search
// object constructors.  NewFromList leaves the result object with zero values
// and does not call the Reset method.
func NewFromList(n int) FromList {
	return FromList{Paths: make([]PathEnd, n)}
}

// Reset initializes a FromList in preparation for a search.  Search methods
// will call this function and you don't typically call it from application
// code.
func (t *FromList) reset() {
	for n := range t.Paths {
		t.Paths[n] = PathEnd{From: -1, Len: 0}
	}
	t.Leaves = big.Int{}
	t.MaxLen = 0
}

// PathTo decodes a FromList, recovering a found path.
//
// The found path is returned as a list of nodes where the first element is
// the start node of the search and the last element is the specified end node.
// If the search did not find a path end the slice result will be nil.
//
// Acceptable end nodes are defined by the specific search.  For most all-paths
// searches, any end node may be specified.  For many single-path searches,
// only the end node specified in the original search is valid.
func (t *FromList) PathTo(end int) []int {
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

// CommonAncestor returns the common ancestor of a and b.
//
// It returns -1 if a or b are invalid node numbers.
func (t *FromList) CommonAncestor(a, b int) int {
	p := t.Paths
	if a < 0 || b < 0 || a >= len(p) || b >= len(p) {
		return -1
	}
	if p[a].Len < p[b].Len {
		a, b = b, a
	}
	for bl := p[b].Len; p[a].Len > bl; {
		a = p[a].From
	}
	for a != b {
		a = p[a].From
		b = p[b].From
	}
	return a
}

// Undirected contructs the undirected graph corresponding to the receiver.
func (t *FromList) Undirected() AdjacencyList {
	g := make(AdjacencyList, len(t.Paths))
	for n, p := range t.Paths {
		if p.From == -1 {
			continue
		}
		g[n] = append(g[n], p.From)
		g[p.From] = append(g[p.From], n)
	}
	return g
}

// Undirected contructs the corresponding undirected graph with edge labels.
func (t *FromList) UndirectedLabeled() LabeledAdjacencyList {
	g := make(LabeledAdjacencyList, len(t.Paths))
	for n, p := range t.Paths {
		if p.From == -1 {
			continue
		}
		g[n] = append(g[n], Half{To: p.From, Label: n})
		g[p.From] = append(g[p.From], Half{n, n})
	}
	return g
}
