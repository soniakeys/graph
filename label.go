// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"math/big"
)

// A LabledAdjacencyList represents a graph as a list of neighbors for each
// node, connected by labeled arcs.
type LabeledAdjacencyList [][]Half

// Half is a half arc, representing a labeled arc and the "neighbor" node
// that the arc leads to.
//
// Halfs can be composed to form a labeled adjacency list.
type Half struct {
	To    int // node ID, usable as a slice index
	Label int // half-arc ID for application data
}

// FromHalf is a half arc, representing a labeled arc and the "neighbor" node
// that the arc originates from.
type FromHalf struct {
	From  int
	Label int
}

// A WeightFunc accesses arc weights from arc labels.
type WeightFunc func(label int) (weight float64)

// NegativeArc returns true if the receiver graph contains a negative arc.
func (g LabeledAdjacencyList) NegativeArc(w WeightFunc) bool {
	for _, nbs := range g {
		for _, nb := range nbs {
			if w(nb.Label) < 0 {
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
func (g LabeledAdjacencyList) ValidTo() bool {
	for _, nbs := range g {
		for _, nb := range nbs {
			if nb.To < 0 || nb.To >= len(g) {
				return false
			}
		}
	}
	return true
}

// LabeledFromTree represents a tree where each node is associated with
// a half arc identifying an arc from another node.
//
// LabeledFromTree is a variant of the FromTree type for labeled graphs.
//
// Paths represents a tree with information about the path to each node
// from a start node.  See LabeledPathEnd documentation.
//
// Leaves is used by some search and traversal functions to return
// extra information.  Where Leaves is used it serves as a bitmap where
// Leave.Bit() == 1 for each leaf of the tree.
//
// Missing, compared to FromTree, is the maximum path length.
//
// In addition to the fields of FromTree, NoPath designates a null distance
// value for non-existant paths and arcs.  You can change this value before
// calling a search or traversal function and it will use NoPath to populate
// distance results where a null value is needed.  The default assigned by
// the constructor is +Inf, which compares well against valid distances.
// Other values you might consider are NaN, which might be considered more
// correct; 0, which sums well with other distances; or -1 which is easily
// tested for as an invalid distance.
//
// Within Paths, LabeledPathEnd.Dist will contain the path distance as the sum
// of arc weights from the start node.  If the node is not reachable, searches
// will set Dist to WeightedFromTree.NoPath.
//
// A single WeightedFromTree can also represent a forest.  In this case paths
// from all leaves do not return to a single start node, but multiple start
// nodes.
type LabeledFromTree struct {
	Paths  []LabeledPathEnd // tree representation
	Leaves big.Int          // leaves of tree
}

// A WeightedPathEnd associates a half arc and a path length.
//
// WeightedPathEnd is an element type of FromTree, a return type from various
// search functions.
//
// For a start node of a search, From.From will be -1 and Len will be 1.
// For other nodes reached by the search, From represents a half arc in
// a path back to start and Len represents the number of nodes in the path.
// For nodes not reached by the search, From.From will be -1 and Len will be 0.
type LabeledPathEnd struct {
	// a "from" half arc, the node the arc comes from and the associated label
	From FromHalf
	Len  int // number of nodes in path from start
}

// NewWeightedFromTree creates a WeightedFromTree object.  You don't typically
// call this function from application code.  Rather it is typically called by
// search object constructors.  NewWeightedFromTree leaves the result object
// with zero values and does not call the Reset method.
func NewLabeledFromTree(n int) *LabeledFromTree {
	t := &LabeledFromTree{
		Paths: make([]LabeledPathEnd, n),
	}
	t.reset()
	return t
}

func (t *LabeledFromTree) reset() {
	for n := range t.Paths {
		t.Paths[n] = LabeledPathEnd{
			From: FromHalf{-1, -1},
		}
	}
	t.Leaves = big.Int{}
}

// PathTo decodes a WeightedFromTree, recovering the found path from start to
// end, where start was an argument to SingleShortestPath or AllShortestPaths
// and end is the argument to this method.
//
// The slice result represents the found path with a sequence of half arcs.
// If no path exists from start to end the slice result will be nil.
// For the first element, representing the start node, the arc weight is
// meaningless and will be WeightedFromTree.NoPath.  The total path distance
// is also returned.  Path distance is the sum of arc weights, excluding of
// couse the meaningless arc weight of the first Half.
func (t *LabeledFromTree) PathTo(end int) []Half {
	n := t.Paths[end].Len
	if n == 0 {
		return nil
	}
	p := make([]Half, n)
	for {
		n--
		f := &t.Paths[end].From
		p[n] = Half{end, f.Label}
		if n == 0 {
			return p
		}
		end = f.From
	}
}

// UnLabeled constructs the equivalent unlabeled graph.
func (g LabeledAdjacencyList) Unlabeled() AdjacencyList {
	a := make(AdjacencyList, len(g))
	for n, nbs := range g {
		to := make([]int, len(nbs))
		for i, nb := range nbs {
			to[i] = nb.To
		}
		a[n] = to
	}
	return a
}
