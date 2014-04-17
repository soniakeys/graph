// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Ed is a simple and fast graph library.
//
// Ed is a graph library of the kind where you create graphs out of
// Ed concrete types, perhaps parallel to existing graph data structures
// in your application.  You call some function such as a graph search
// on the Ed graph, then use the result to navigate your application data.
//
// Ed graphs contain only data minimally neccessary for search functions.
// This minimalism simplifies Ed code and allows faster searches.  Zero-based
// integer node IDs serve directly as slice indexes.  Nodes and edge objects
// are structs rather than interfaces.  Maps are not needed to associate
// arbitrary IDs with node or edge types.  Ed graphs are memory efficient
// and large graphs can potentially be handled, especially if Ed graphs are
// constructed in an online manner.
//
// Terminology
//
// There are lots of terms to pick from.  Goals for picking terms for this
// this package include picking popular terms, terms that reduce confusion,
// terms that describe, and terms that read well.
//
// This package uses the term "node" rather than "vertex."  It uses "arc"
// to mean a directed edge, and uses "from" and "to" to refer to the ends
// of an arc.  It uses "start" and "end" to refer to endpoints of a search
// or traversal.
//
// A float64 value associated with an arc is "weight."  The sum of arc weights
// along a path is a "distance."  The number of nodes in a path is the path's
// "length."
//
// A "half arc" represents just one end of an arc, perhaps assocating it with
// an arc weight.  The more common half to work with is the "to half" (the
// type name is simply "Half".)  A list of half arcs can represent a
// "neighbor list," neighbors of a single node.  A list of neighbor lists
// forms an "adjacency list" which represents a directed graph.
//
// A node that is a neighbor of itself represents a "loop."  Duplicate
// neighbors (when a node appears more than once in the same neighbor list)
// represent "parallel arcs."
//
// Finally, this package documentation takes back the word "object" to
// refer to a Go value, especially a value of a type with methods.
package ed

import (
	"math/big"
	"sort"
)

// file ed.go contains definitions for unweighted graphs

// An AdjacencyList represents a graph as a list of neighbors for each node.
// The "node ID" of a node is simply it's slice index in the AdjacencyList.
//
// Adjacency lists are inherently directed. To represent an undirected graph,
// create reciprocal neighbors.
type AdjacencyList [][]int

// Valid validates that no arcs in the reciever graph lead outside the graph.
//
// Ints in an adjacency list structure represent half arcs.  Valid
// returns true if all int values are valid slice indexes back into g.
func (g AdjacencyList) Valid() bool {
	for _, nbs := range g {
		for _, nb := range nbs {
			if nb < 0 || nb >= len(g) {
				return false
			}
		}
	}
	return true
}

// Undirected determines if an adjacency list is undirected.
//
// Undirected validates that g is undirected by checking that every neighbor
// has a reciprocal.  That is, that for every arc from->to, that to->from also
// exists.
//
// If the graph is undirected, Undirected returns true, -1, -1.  If an arc
// is found witout a reciprocal, Undirected returns false along with from
// and to nodes of the arc.
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

// Simple checks for loops and parallel arcs.
//
// Simple returns true, -1 for simple graphs.  If a loop or parallel arc is found,
// simple returns false and a node that has a loop or parallel arc in its neighbor list.
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

// Transpose constructs a new adjacency list that is the transpose of g.
//
// For every arc from->to of g, the result will have an arc to->from.
// Transpose also returns m the number of arcs in g (equal to the number of
// arcs in the result.)
func (g AdjacencyList) Transpose() (t AdjacencyList, m int) {
	t = make([][]int, len(g))
	for n, nbs := range g {
		for _, nb := range nbs {
			t[nb] = append(t[nb], n)
			m++
		}
	}
	return
}

// ConnectedComponents determines the connected components in an undirected
// graph.
//
// Returned is a slice with a single representative node from each connected
// component.
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

// Bipartite determines if a graph is bipartite.
//
// If so, Bipartite returns true and the two-coloring of the graph.  Each
// color set is returned as a bitmap.  If the graph is not bipartite,
// Bipartite returns false and an odd cycle as an int slice.
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

// Acyclic determines if a directed graph contains cycles.
//
// Acyclic returns true if there are no cycles.
// Acyclic returns false if a cycle is detected.
func (g AdjacencyList) Acyclic() bool {
	a := true
	var temp, perm big.Int
	var df func(int)
	df = func(n int) {
		switch {
		case temp.Bit(n) == 1:
			a = false
			return
		case perm.Bit(n) == 1:
			return
		}
		temp.SetBit(&temp, n, 1)
		for _, nb := range g[n] {
			df(nb)
			if !a {
				return
			}
		}
		temp.SetBit(&temp, n, 0)
		perm.SetBit(&perm, n, 1)
	}
	for n := range g {
		if perm.Bit(n) == 1 {
			continue
		}
		df(n)
		if !a {
			return false
		}
	}
	return true
}

// Topological computes a topological sort of a directed acyclic graph.
//
// Returned is a permutation of node numbers in topologically sorted order.
// If the graph is found not to be acyclic, Topological returns nil.
func (g AdjacencyList) Topological() []int {
	t := make([]int, len(g))
	i := len(t) - 1
	a := true
	var temp, perm big.Int
	var df func(int)
	df = func(n int) {
		switch {
		case temp.Bit(n) == 1:
			a = false
			return
		case perm.Bit(n) == 1:
			return
		}
		temp.SetBit(&temp, n, 1)
		for _, nb := range g[n] {
			df(nb)
			if !a {
				return
			}
		}
		temp.SetBit(&temp, n, 0)
		perm.SetBit(&perm, n, 1)
		t[i] = n
		i--
	}
	for n := range g {
		if perm.Bit(n) == 1 {
			continue
		}
		df(n)
		if !a {
			return nil
		}
	}
	return t
}

// Tarjan identifies strongly connected components in a directed graph using
// Tarjan's algorithm.
//
// Returned is a list of components, each component is a list of nodes.
func (g AdjacencyList) Tarjan() (scc [][]int) {
	// straight from WP
	var indexed, stacked big.Int
	index := make([]int, len(g))
	lowlink := make([]int, len(g))
	x := 0
	var S []int
	var sc func(int)
	sc = func(n int) {
		index[n] = x
		indexed.SetBit(&indexed, n, 1)
		lowlink[n] = x
		x++
		S = append(S, n)
		stacked.SetBit(&stacked, n, 1)
		for _, nb := range g[n] {
			if indexed.Bit(nb) == 0 {
				sc(nb)
				if lowlink[nb] < lowlink[n] {
					lowlink[n] = lowlink[nb]
				}
			} else if stacked.Bit(nb) == 1 {
				if index[nb] < lowlink[n] {
					lowlink[n] = index[nb]
				}
			}
		}
		if lowlink[n] == index[n] {
			var c []int
			for {
				last := len(S) - 1
				w := S[last]
				S = S[:last]
				stacked.SetBit(&stacked, w, 0)
				c = append(c, w)
				if w == n {
					scc = append(scc, c)
					break
				}
			}
		}
	}
	for n := range g {
		if indexed.Bit(n) == 0 {
			sc(n)
		}
	}
	return scc
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
// search to the given end node.
//
// The found path is returned as a list of nodes.
// If the search did not find a path end the slice result will be nil.
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
