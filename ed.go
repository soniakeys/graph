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
// an arc weight.  The more common half to work with is the "half to" (the
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
	"math"
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

// A WeightedAdjacencyList represents a graph as a list of neighbors for each
// node, connected by weighted arcs.
type WeightedAdjacencyList [][]Half

// Half is a half arc, representing a weighted arc and the "neighbor" node
// that the arc leads to.
//
// Halfs can be composed to form a weighted adjacency list.
type Half struct {
	To        int // node ID, usable as a slice index
	ArcWeight float64
}

// HalfFrom is a half arc, representing a weighted arc and the "neighbor" node
// that the arc originates from.
type HalfFrom struct {
	From      int
	ArcWeight float64
}

// NegativeArc returns true if the receiver graph contains a negative arc.
func (a WeightedAdjacencyList) NegativeArc() bool {
	for _, nbs := range a {
		for _, nb := range nbs {
			if nb.ArcWeight < 0 {
				return true
			}
		}
	}
	return false
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

// ValidTo validates that no arcs in the reciever graph lead outside the graph.
//
// ValidTo returns true if all Half.To values are valid slice indexes
// back into g.
func (g WeightedAdjacencyList) ValidTo() bool {
	for _, nbs := range g {
		for _, nb := range nbs {
			if nb.To < 0 || nb.To >= len(g) {
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

// WeightedFromTree is a variant of the FromTree type for weighted graphs.
//
// In addition to the fields of FromTree, NoPath designates a null distance
// value for non-existant paths and arcs.  You can change this value before
// calling a search or traversal function and it will use NoPath to populate
// distance results where a null value is needed.  The default assigned by
// the constructor is -Inf, which compares well against valid distances.
// Other values you might consider are NaN, which might be considered more
// correct; 0, which sums well with other distances; or -1 which is easily
// tested for as an invalid distance.
//
// Within Paths, WeightedPathDist will contain the path distance as the sum
// of arc weights from the start node.  If the node is not reachable, searches
// will set PathDist to WeightedFromTree.NoPath.
//
// Missing, compared to FromTree is the maximum path length.
type WeightedFromTree struct {
	Start  int               // start node, argument to the search function
	Paths  []WeightedPathEnd // tree representation
	NoPath float64           // null value
}

// A WeightedPathEnd associates a half arc and a path length.
//
// See WeightedFromTree for use by search functions.
type WeightedPathEnd struct {
	From HalfFrom // half arc
	Dist float64  // path distance, sum of arc weights from start
	Len  int      // number of nodes in path from start
}

// NewWeightedFromTree creates a WeightedFromTree object.  You don't typically
// call this function from application code.  Rather it is typically called by
// search object constructors.  NewWeightedFromTree leaves the result object
// with zero values and does not call the Reset method.
func NewWeightedFromTree(n int) *WeightedFromTree {
	return &WeightedFromTree{
		Paths:  make([]WeightedPathEnd, n),
		NoPath: math.Inf(1),
	}
}

// Reset initializes a WeightedFromTree in preparation for a search.
// Search methods will call this function and you don't typically call it
// from application code.
func (t *WeightedFromTree) Reset() {
	t.Start = -1
	for n := range t.Paths {
		t.Paths[n] = WeightedPathEnd{
			Dist: t.NoPath,
			From: HalfFrom{-1, t.NoPath},
		}
	}
}

// PathTo decodes Dijkstra.Result, recovering the found path from start to
// end, where start was an argument to SingleShortestPath or AllShortestPaths
// and end is the argument to this method.
//
// The slice result represents the found path with a sequence of half arcs.
// If no path exists from start to end the slice result will be nil.
// For the first element, representing the start node, the arc weight is
// meaningless and will be Dijkstra.NoPath.  The total path distance is also
// returned.  Path distance is the sum of arc weights, excluding of couse
// the meaningless arc weight of the first Half.
func (t *WeightedFromTree) PathTo(end int) ([]Half, float64) {
	n := t.Paths[end].Len
	if n == 0 {
		return nil, t.NoPath
	}
	p := make([]Half, n)
	dist := 0.
	for {
		n--
		f := &t.Paths[end].From
		p[n] = Half{end, f.ArcWeight}
		if n == 0 {
			return p, dist
		}
		dist += f.ArcWeight
		end = f.From
	}
}
