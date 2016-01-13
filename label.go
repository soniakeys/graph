// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"math"
	"sort"
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

// WeightFunc returns a weight for a given label.
//
// WeightFunc is a parameter type for various search functions.  The intent
// is to return a weight corresponding to an arc label.  The name "weight"
// is an abstract term.  An arc "weight" will typically have some application
// specific meaning other than physical weight.
type WeightFunc func(label int) (weight float64)

// AddEdge adds an edge to a labeled graph.
//
// It can be useful for constructing undirected graphs.
//
// When n1 and n2 are distinct, it adds the arc n1->n2 and the reciprocal
// n2->n1.  When n1 and n2 are the same, it adds a single arc loop.
//
// If the edge already exists in *p, a parallel edge is added.
//
// The pointer receiver allows the method to expand the graph as needed
// to include the values n1 and n2.  If n1 or n2 happen to be greater than
// len(*p) the method does not panic, but simply expands the graph.
func (p *LabeledAdjacencyList) AddEdge(e LabeledEdge) {
	// Similar code in AdjacencyList.AddEdge.

	// determine max of the two end points
	max := e.N1
	if e.N2 > max {
		max = e.N2
	}
	// expand graph if needed, to include both
	g := *p
	if max >= len(g) {
		*p = make(LabeledAdjacencyList, max+1)
		copy(*p, g)
		g = *p
	}
	// create one half-arc,
	g[e.N1] = append(g[e.N1], Half{To: e.N2, Label: e.Label})
	// and except for loops, create the reciprocal
	if e.N1 != e.N2 {
		g[e.N2] = append(g[e.N2], Half{To: e.N1, Label: e.Label})
	}
}

// DAGMaxLenPath finds a maximum length path in a directed acyclic graph.
//
// Length here means number of nodes or arcs, not a sum of arc weights.
//
// Argument ordering must be a topological ordering of g.
//
// Returned is a node beginning a maximum length path, and a path of arcs
// starting from that node.
func (g LabeledAdjacencyList) DAGMaxLenPath(ordering []int) (n int, path []Half) {
	// dynamic programming. visit nodes in reverse order. for each, compute
	// longest path as one plus longest of 'to' nodes.
	// Visits each arc once.  Time complexity O(m).
	//
	// Similar code in dir.go.
	mlp := make([][]Half, len(g)) // index by node number
	for i := len(ordering) - 1; i >= 0; i-- {
		fr := ordering[i] // node number
		to := g[fr]
		if len(to) == 0 {
			continue
		}
		mt := to[0]
		for _, to := range to[1:] {
			if len(mlp[to.To]) > len(mlp[mt.To]) {
				mt = to
			}
		}
		p := append([]Half{mt}, mlp[mt.To]...)
		mlp[fr] = p
		if len(p) > len(path) {
			n = fr
			path = p
		}
	}
	return
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

// IsUndirected returns true if g represents an undirected graph.
//
// Returns true when all non-loop arcs are paired in reciprocal pairs with
// matching labels.  Otherwise returns false and an example unpaired arc.
//
// Note the requirement that reciprocal pairs have matching labels is
// an additional test not present in the otherwise equivalent unlabled version
// of IsUndirected.
func (g LabeledAdjacencyList) IsUndirected() (u bool, from int, to Half) {
	unpaired := make(LabeledAdjacencyList, len(g))
	for fr, to := range g {
	arc: // for each arc in g
		for _, to := range to {
			if to.To == fr {
				continue // loop
			}
			// search unpaired arcs
			ut := unpaired[to.To]
			for i, u := range ut {
				if u.To == fr && u.Label == to.Label { // found reciprocal
					last := len(ut) - 1
					ut[i] = ut[last]
					unpaired[to.To] = ut[:last]
					continue arc
				}
			}
			// reciprocal not found
			unpaired[fr] = append(unpaired[fr], to)
		}
	}
	for fr, to := range unpaired {
		if len(to) > 0 {
			return false, fr, to[0]
		}
	}
	return true, -1, to
}

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

// Simple checks for loops and parallel arcs.
//
// A graph is "simple" if it has no loops or parallel arcs.
//
// Simple returns true, -1 for simple graphs.  If a loop or parallel arc is
// found, simple returns false and and a node that represents a counterexample
// to the graph being simple.
//
// See also the eqivalent unlabeled Simple.
func (g LabeledAdjacencyList) Simple() (s bool, n int) {
	var t []int
	for n, nbs := range g {
		if len(nbs) == 0 {
			continue
		}
		t = t[:0]
		for _, nb := range nbs {
			t = append(t, nb.To)
		}
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

// TarjanBiconnectedComponents, for undirected simple graphs.
//
// A list of components is returned, with each component represented as an
// edge list.
//
// See also the eqivalent unlabeled TarjanBiconnectedComponents.
func (g LabeledAdjacencyList) TarjanBiconnectedComponents() (components [][]LabeledEdge) {
	// Implemented closely to pseudocode in "Depth-first search and linear
	// graph algorithms", Robert Tarjan, SIAM J. Comput. Vol. 1, No. 2,
	// June 1972.
	//
	// Note Tarjan's "adjacency structure" is graph.AdjacencyList,
	// His "adjacency list" is an element of a graph.AdjacencyList, also
	// termed a "to-list", "neighbor list", or "child list."
	//
	// Nearly identical code in undir.go.
	number := make([]int, len(g))
	lowpt := make([]int, len(g))
	var stack []LabeledEdge
	var i int
	var biconnect func(int, int)
	biconnect = func(v, u int) {
		i++
		number[v] = i
		lowpt[v] = i
		for _, w := range g[v] {
			if number[w.To] == 0 {
				stack = append(stack, LabeledEdge{Edge{v, w.To}, w.Label})
				biconnect(w.To, v)
				if lowpt[w.To] < lowpt[v] {
					lowpt[v] = lowpt[w.To]
				}
				if lowpt[w.To] >= number[v] {
					var bcc []LabeledEdge
					top := len(stack) - 1
					for number[stack[top].N1] >= number[w.To] {
						bcc = append(bcc, stack[top])
						stack = stack[:top]
						top--
					}
					bcc = append(bcc, stack[top])
					stack = stack[:top]
					top--
					components = append(components, bcc)
				}
			} else if number[w.To] < number[v] && w.To != u {
				stack = append(stack, LabeledEdge{Edge{v, w.To}, w.Label})
				if number[w.To] < lowpt[v] {
					lowpt[v] = number[w.To]
				}
			}
		}
	}
	for w := range g {
		if number[w] == 0 {
			biconnect(w, 0)
		}
	}
	return
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

// Undirected returns copy of g augmented as needed to make it undirected,
// with reciprocal arcs having matching labels.
func (g LabeledAdjacencyList) UndirectedCopy() LabeledAdjacencyList {
	c, _ := g.Copy()                         // start with a copy
	rw := make(LabeledAdjacencyList, len(g)) // "reciprocals wanted"
	for fr, to := range g {
	arc: // for each arc in g
		for _, to := range to {
			if to.To == fr {
				continue // arc is a loop
			}
			// search wanted arcs
			wf := rw[fr]
			for i, w := range wf {
				if w == to { // found, remove
					last := len(wf) - 1
					wf[i] = wf[last]
					rw[fr] = wf[:last]
					continue arc
				}
			}
			// arc not found, add to reciprocal to wanted list
			rw[to.To] = append(rw[to.To], Half{To: fr, Label: to.Label})
		}
	}
	// add missing reciprocals
	for fr, to := range rw {
		c[fr] = append(c[fr], to...)
	}
	return c
}

// Unlabeled constructs the unlabeled graph corresponding to g.
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

// UnlabeledTranspose, for directed graphs, constructs a new adjacency list
// that is the unlabeled transpose of g.
//
// For every arc from->to of g, the result will have an arc to->from.
// Transpose also counts arcs as it traverses and returns m the number of arcs
// in g (equal to the number of arcs in the result.)
//
// It is equivalent to g.Unlabeled().Transpose() but constructs the result
// directly.
func (g LabeledAdjacencyList) UnlabeledTranspose() (t AdjacencyList, m int) {
	t = make(AdjacencyList, len(g))
	for n, nbs := range g {
		for _, nb := range nbs {
			t[nb.To] = append(t[nb.To], n)
			m++
		}
	}
	return
}

// LabeledEdge is an undirected edge with an associated label.
type LabeledEdge struct {
	Edge
	Label int
}

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

// Len implements sort.Interface.
func (l WeightedEdgeList) Len() int { return len(l.Edges) }

// Less implements sort.Interface.
func (l WeightedEdgeList) Less(i, j int) bool {
	return l.WeightFunc(l.Edges[i].Label) < l.WeightFunc(l.Edges[j].Label)
}

// Swap implements sort.Interface.
func (l WeightedEdgeList) Swap(i, j int) {
	l.Edges[i], l.Edges[j] = l.Edges[j], l.Edges[i]
}
