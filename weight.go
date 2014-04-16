// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"math"
)

// file weight.go contains definitions for weighted graphs

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

// FromHalf is a half arc, representing a weighted arc and the "neighbor" node
// that the arc originates from.
type FromHalf struct {
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
	From FromHalf // half arc
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

func (t *WeightedFromTree) reset() {
	t.Start = -1
	for n := range t.Paths {
		t.Paths[n] = WeightedPathEnd{
			Dist: t.NoPath,
			From: FromHalf{-1, t.NoPath},
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
