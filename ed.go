// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Ed is a simple and fast graph library.
//
// For goals of speed and simplicity, Ed uses zero-based integer node IDs
// and omits interfaces that would accomodate user data or user implemented
// behavior.
//
// To use Ed functions, you typically create a data structure parallel
// to your application data, call an Ed function, and use the result to
// navigate your application data.
//
// Terminology
//
// This package uses the term "node" rather than "vertex."  It uses "arc"
// to mean a directed edge, and uses "from" and "to" to refer to the ends
// of an arc.
//
// A node that is a neighbor of itself is termed a "loop."  Duplicate
// neighbors, or when a node appears more than once in the same neighbor list
// are termed "parallel arcs."
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

// DijkstraResult contains search results for a single node.  A Dijkstra
// object contains a DijkstraResult slice parallel to the input graph.
//
// Following an AllShortestPaths search, PathDist will contain the path
// distance as the sum of arc weights for the found shortest path from the
// start node to the node corresponding to the Dijkstra result.  If the node
// is not reachable, PathDist will be Dijkstra.NoPath.
//
// The FromTree fields for all reachable nodes represent a spanning
// tree encoding all shortest paths.  The structure is an "inverse
// arborescence", "in-tree", or "spaghetti stack."  For unreachable nodes
// FromTree.From will be -1, FromTree.ArcWeight will be Dijkstra.NoPath.
// See Dijkstra.PathTo for convenient decoding of this structure.
//
// PathLen will contain the number of nodes in the found shortest path, or
// 0 if the node is unreachable.  Testing PathLen == 0 is the simplest test
// of whether a node is reachable.  Note that in the case where there are
// multiple paths with same minimum PathDist, PathLen contains only the length
// of the path encoded by FromTree.  Other paths of minimum distance may have
// fewer nodes.
//
// Following a SingleShortestPath search, results are only valid along the
// found path.  Results for other nodes are meaningless.
type FromTree struct {
	Start  int
	Paths  []PathEnd
	MaxLen int
}

func NewFromTree(n int) *FromTree {
	return &FromTree{Paths: make([]PathEnd, n)}
}

func (t *FromTree) Reset() {
	t.Start = -1
	for n := range t.Paths {
		t.Paths[n] = PathEnd{From: -1, Len: 0}
	}
	t.MaxLen = 0
}

type PathEnd struct {
	From int
	Len  int
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

type WeightedFromTree struct {
	Start  int
	Paths  []WeightedPathEnd
	NoPath float64
}

type WeightedPathEnd struct {
	Dist float64
	From HalfFrom
	Len  int
}

func NewWeightedFromTree(n int) *WeightedFromTree {
	return &WeightedFromTree{
		Paths:  make([]WeightedPathEnd, n),
		NoPath: math.Inf(1),
	}
}

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
