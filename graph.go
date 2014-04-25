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

// A FromTree represents a tree where each node is associated with
// a half arc identifying an arc from another node.
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
func newFromTree(n int) *FromTree {
	return &FromTree{Paths: make([]PathEnd, n)}
}

// Reset initializes a FromTree in preparation for a search.  Search methods
// will call this function and you don't typically call it from application
// code.
func (t *FromTree) reset() {
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
// The found path is returned as a list of nodes where the first element is
// the start node of the search and the last element is the specified end node.
// If the search did not find a path end the slice result will be nil.
//
// Acceptable end nodes are defined by the specific search.  For most all-paths
// searches, any end node may be specified.  For many single-path searches,
// only the end node specified in the original search is valid.
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
