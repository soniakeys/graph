// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import "math"

// graph.go contains type definitions for all graph types and components.
// Also, go generate directives for source transformations.
//
// For readability, the types are defined in a dependency order:
//
//  NI
//  AdjacencyList
//  Directed
//  Undirected
//  LI
//  Half
//  LabeledAdjacencyList
//  LabeledDirected
//  LabeledUndirected
//  Edge
//  LabeledEdge
//  WeightFunc
//  WeightedEdgeList
//  TraverseOption

//go:generate cp adj_cg.go adj_RO.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w adj_RO.go
//go:generate gofmt -r "n.To -> n" -w adj_RO.go
//go:generate gofmt -r "Half -> NI" -w adj_RO.go

//go:generate cp dir_cg.go dir_RO.go
//go:generate gofmt -r "LabeledDirected -> Directed" -w dir_RO.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w dir_RO.go
//go:generate gofmt -r "labEulerian -> eulerian" -w dir_RO.go
//go:generate gofmt -r "newLabEulerian -> newEulerian" -w dir_RO.go
//go:generate gofmt -r "Half{n, -1} -> n" -w dir_RO.go
//go:generate gofmt -r "n.To -> n" -w dir_RO.go
//go:generate gofmt -r "Half -> NI" -w dir_RO.go

//go:generate cp undir_cg.go undir_RO.go
//go:generate gofmt -r "LabeledUndirected -> Undirected" -w undir_RO.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w undir_RO.go
//go:generate gofmt -r "newLabEulerian -> newEulerian" -w undir_RO.go
//go:generate gofmt -r "Half{n, -1} -> n" -w undir_RO.go
//go:generate gofmt -r "n.To -> n" -w undir_RO.go
//go:generate gofmt -r "Half -> NI" -w undir_RO.go

// An AdjacencyList represents a graph as a list of neighbors for each node.
// The "node ID" of a node is simply it's slice index in the AdjacencyList.
// For an AdjacencyList g, g[n] represents arcs going from node n to nodes
// g[n].
//
// Adjacency lists are inherently directed but can be used to represent
// directed or undirected graphs.  See types Directed and Undirected.
type AdjacencyList [][]NI

// Directed represents a directed graph.
//
// Directed methods generally rely on the graph being directed, specifically
// that arcs do not have reciprocals.
type Directed struct {
	AdjacencyList // embedded to include AdjacencyList methods
}

// Undirected represents an undirected graph.
//
// In an undirected graph, for each arc between distinct nodes there is also
// a reciprocal arc, an arc in the opposite direction.  Loops do not have
// reciprocals.
//
// Undirected methods generally rely on the graph being undirected,
// specifically that every arc between distinct nodes has a reciprocal.
type Undirected struct {
	AdjacencyList // embedded to include AdjacencyList methods
}

// LI is a label integer, used for associating labels with arcs.
type LI int32

// Half is a half arc, representing a labeled arc and the "neighbor" node
// that the arc leads to.
//
// Halfs can be composed to form a labeled adjacency list.
type Half struct {
	To    NI // node ID, usable as a slice index
	Label LI // half-arc ID for application data, often a weight
}

// A LabeledAdjacencyList represents a graph as a list of neighbors for each
// node, connected by labeled arcs.
//
// Arc labels are not necessarily unique arc IDs.  Different arcs can have
// the same label.
//
// Arc labels are commonly used to assocate a weight with an arc.  Arc labels
// are general purpose however and can be used to associate arbitrary
// information with an arc.
//
// Methods implementing weighted graph algorithms will commonly take a
// weight function that turns a label int into a float64 weight.
//
// If only a small amount of information -- such as an integer weight or
// a single printable character -- needs to be associated, it can sometimes
// be possible to encode the information directly into the label int.  For
// more generality, some lookup scheme will be needed.
//
// In an undirected labeled graph, reciprocal arcs must have identical labels.
// Note this does not preclude parallel arcs with different labels.
type LabeledAdjacencyList [][]Half

// LabeledDirected represents a directed labeled graph.
//
// This is the labeled version of Directed.  See types LabeledAdjacencyList
// and Directed.
type LabeledDirected struct {
	LabeledAdjacencyList // embedded to include LabeledAdjacencyList methods
}

// LabeledUndirected represents an undirected labeled graph.
//
// This is the labeled version of Undirected.  See types LabeledAdjacencyList
// and Undirected.
type LabeledUndirected struct {
	LabeledAdjacencyList // embedded to include LabeledAdjacencyList methods
}

// Edge is an undirected edge between nodes N1 and N2.
type Edge struct{ N1, N2 NI }

// LabeledEdge is an undirected edge with an associated label.
type LabeledEdge struct {
	Edge
	LI
}

// WeightFunc returns a weight for a given label.
//
// WeightFunc is a parameter type for various search functions.  The intent
// is to return a weight corresponding to an arc label.  The name "weight"
// is an abstract term.  An arc "weight" will typically have some application
// specific meaning other than physical weight.
type WeightFunc func(label LI) (weight float64)

// WeightedEdgeList is a graph representation.
//
// It is a labeled edge list, with an associated weight function to return
// a weight given an edge label.
//
// Also associated is the order, or number of nodes of the graph.
// All nodes occurring in the edge list must be strictly less than Order.
//
// WeigtedEdgeList sorts by weight, obtained by calling the weight function.
// If weight computation is expensive, consider supplying a cached or
// memoized version.
type WeightedEdgeList struct {
	Order int
	WeightFunc
	Edges []LabeledEdge
}

// DistanceMatrix constructs a distance matrix corresponding to the weighted
// edges of l.
//
// An edge n1, n2 with WeightFunc return w is represented by both
// d[n1][n2] == w and d[n2][n1] = w.  In case of parallel edges, the lowest
// weight is stored.  The distance from any node to itself d[n][n] is 0, unless
// the node has a loop with a negative weight.  If g has no edge between n1 and
// distinct n2, +Inf is stored for d[n1][n2] and d[n2][n1].
//
// The returned DistanceMatrix is suitable for DistanceMatrix.FloydWarshall.
func (l WeightedEdgeList) DistanceMatrix() (d DistanceMatrix) {
	d = newDM(l.Order)
	for _, e := range l.Edges {
		n1 := e.Edge.N1
		n2 := e.Edge.N2
		wt := l.WeightFunc(e.LI)
		// < to pick min of parallel arcs (also nicely ignores NaN)
		if wt < d[n1][n2] {
			d[n1][n2] = wt
			d[n2][n1] = wt
		}
	}
	return
}

// A DistanceMatrix is a square matrix representing some distance between
// nodes of a graph.  If the graph is directected, d[from][to] represents
// some distance from node 'from' to node 'to'.  Depending on context, the
// distance may be an arc weight or path distance.  A value of +Inf typically
// means no arc or no path between the nodes.
type DistanceMatrix [][]float64

// little helper function, makes a blank distance matrix for FloydWarshall.
// could be exported?
func newDM(n int) DistanceMatrix {
	inf := math.Inf(1)
	d := make(DistanceMatrix, n)
	for i := range d {
		di := make([]float64, n)
		for j := range di {
			di[j] = inf
		}
		di[i] = 0
		d[i] = di
	}
	return d
}

// FloydWarshall finds all pairs shortest distances for a weighted graph
// without negative cycles.
//
// It operates on a distance matrix representing arcs of a graph and
// destructively replaces arc weights with shortest path distances.
//
// In result array d, d[fr][to] will be the shortest distance from node 'fr'
// to node 'to'.  An element value of +Inf means no path exists.  Any diagonal
// element < 0 indicates a negative cycle exists.
//
// See DistanceMatrix constructor methods of LabeledAdjacencyList and
// WeightedEdgeList for suitable inputs.
func (d DistanceMatrix) FloydWarshall() {
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
