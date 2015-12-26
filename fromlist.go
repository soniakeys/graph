// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// fromlist.go

package graph

import "math/big"

// FromList represents a tree where each node is associated with
// a half arc identifying an arc from another node.
//
// The Paths member represents the tree structure.  Leaves and MaxLen are
// not always needed.  Where Leaves is used it serves as a bitmap where
// Leave.Bit() == 1 for each leaf of the tree.  Where MaxLen is used it is
// provided primarily as a convenience for functions that might want to
// anticipate the maximum path length that would be encountered traversing
// the tree.
//
// Various graph search functions use a FromList to returns search results.
// For a start node of a search, From will be -1 and Len will be 1. For other
// nodes reached by the search, From represents a half arc in a path back to
// start and Len represents the number of nodes in the path.  For nodes not
// reached by the search, From will be -1 and Len will be 0.
//
// A single FromList can also represent a forest.  In this case paths from
// all leaves do not return to a single root node, but multiple root nodes.
type FromList struct {
	Paths  []PathEnd // tree representation
	Leaves big.Int   // leaves of tree
	MaxLen int       // length of longest path, max of all PathEnd.Len values
}

// PathEnd associates a half arc and a path length.
//
// A PathEnd list is an element type of FromList.
type PathEnd struct {
	From int // a "from" half arc, the node the arc comes from
	Len  int // number of nodes in path from start
}

// NewFromList creates a FromList object of given order.
func NewFromList(n int) FromList {
	return FromList{Paths: make([]PathEnd, n)}
}

// reset initializes a FromList in preparation for a search.  Search methods
// will call this function and you don't typically call it from application
// code.
func (t *FromList) reset() {
	for n := range t.Paths {
		t.Paths[n] = PathEnd{From: -1, Len: 0}
	}
	t.Leaves = big.Int{}
	t.MaxLen = 0
}

// BoundsOk validates the "from" values in the list.
//
// Negative values are allowed as they indicate root nodes.
//
// BoundsOk returns true when all from values are less than len(t).
// Otherwise it returns false and a node with a from value >= len(t).
func (t *FromList) BoundsOk() (ok bool, n int) {
	for n, e := range t.Paths {
		if e.From >= len(t.Paths) {
			return false, n
		}
	}
	return true, -1
}

// PathTo decodes a FromList, recovering a single path.
//
// The path is returned as a list of nodes where the first element will be
// a root node and the last element will be the specified end node.
//
// Only the Paths member of the receiver is used.  Other members of the
// FromList do not need to be valid, however the MaxLen member can be useful
// for allocating argument p.
//
// Argument p can provide the result slice.  If p has capacity for the result
// it will be used, otherwise a new slice is created for the result.
//
// See also function PathTo.
func (t *FromList) PathTo(end int, p []int) []int {
	return PathTo(t.Paths, end, p)
}

// PathTo decodes a single path from a PathEnd list.
//
// The path is returned as a list of nodes where the first element will be
// a root node and the last element will be the specified end node.
//
// Argument p can provide the result slice.  If p has capacity for the result
// it will be used, otherwise a new slice is created for the result.
//
// See also method FromList.PathTo.
func PathTo(paths []PathEnd, end int, p []int) []int {
	n := paths[end].Len
	if n == 0 {
		return nil
	}
	if cap(p) >= n {
		p = p[:n]
	} else {
		p = make([]int, n)
	}
	for {
		n--
		p[n] = end
		if n == 0 {
			return p
		}
		end = paths[end].From
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
