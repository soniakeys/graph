// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// fromlist.go

package graph

import (
	"math/big"
)

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

// Transpose contructs the directed graph with arcs in the opposite direction
// of the FromList.  That is, from root toward leaves.
func (t *FromList) Transpose() AdjacencyList {
	g := make(AdjacencyList, len(t.Paths))
	for n, p := range t.Paths {
		if p.From == -1 {
			continue
		}
		g[p.From] = append(g[p.From], n)
	}
	return g
}

// TransposeLabled contructs the labeled directed graph with arcs in the
// opposite direction of the FromList.  That is, from root toward leaves.
func (t *FromList) TransposeLabeled() LabeledAdjacencyList {
	g := make(LabeledAdjacencyList, len(t.Paths))
	for n, p := range t.Paths {
		if p.From == -1 {
			continue
		}
		g[p.From] = append(g[p.From], Half{n, n})
	}
	return g
}

// Undirected contructs the undirected graph corresponding to the FromList.
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

// UndirectedLabeled contructs the corresponding undirected graph with
// edge labels.
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
