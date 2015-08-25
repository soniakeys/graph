// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"math"
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

// Transpose, for directed graphs, constructs a new adjacency list that is
// the transpose of g.
//
// For every arc from->to of g, the result will have an arc to->from.
// Transpose also counts arcs as it traverses and returns m the number of arcs
// in g (equal to the number of arcs in the result.)
func (g LabeledAdjacencyList) Transpose() (t LabeledAdjacencyList, m int) {
	t = make(LabeledAdjacencyList, len(g))
	for n, nbs := range g {
		for _, nb := range nbs {
			t[nb.To] = append(t[nb.To], Half{To: n, Label: nb.Label})
			m++
		}
	}
	return
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
// A single LabeledFromTree can also represent a forest.  In this case paths
// from all leaves do not return to a single start node, but multiple start
// nodes.
type LabeledFromTree struct {
	Paths  []LabeledPathEnd // tree representation
	Leaves big.Int          // leaves of tree
}

// A LabeledPathEnd associates a half arc and a path length.
//
// LabeledPathEnd is an element type of FromTree, a return type from various
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

// NewLabeledFromTree creates a LabeledFromTree object.  You don't typically
// call this function from application code.  Rather it is typically called by
// search object constructors.  NewLabeledFromTree leaves the result object
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

// PathTo decodes a LabeledFromTree, recovering the found path from start to
// end, where start was an argument to SingleShortestPath or AllShortestPaths
// and end is the argument to this method.
//
// The slice result represents the found path with a sequence of half arcs.
// If no path exists from start to end the slice result will be nil.
// For the first element, representing the start node, the arc label is
// meaningless and will be -1.
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

// FloydWarshall finds all pairs shortest distances for a simple weighted
// graph without negative cycles.
//
// In result array d, d[i][j] will be the shortest distance from node i
// to node j.  Any diagonal element < 0 indicates a negative cycle exists.
//
// If g is an undirected graph with no negative edge weights, the result
// array will be a distance matrix, for example as used by package
// github.com/soniakeys/cluster.
func (g LabeledAdjacencyList) FloydWarshall(w WeightFunc) (d [][]float64) {
	d = newFWd(len(g))
	for fr, to := range g {
		for _, to := range to {
			d[fr][to.To] = w(to.Label)
		}
	}
	solveFW(d)
	return
}

// little helper function, makes a blank matrix for FloydWarshall.
func newFWd(n int) [][]float64 {
	d := make([][]float64, n)
	for i := range d {
		di := make([]float64, n)
		for j := range di {
			if j != i {
				di[j] = math.Inf(1)
			}
		}
		d[i] = di
	}
	return d
}

// Floyd Warshall solver, once the matrix d is initialized by arc weights.
func solveFW(d [][]float64) {
	for k, dk := range d {
		for _, di := range d {
			dik := di[k]
			for j := range d {
				if d2 := dik + dk[j]; d2 < di[j] {
					di[j] = d2
				}
			}
		}
	}
}
