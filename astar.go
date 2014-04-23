// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"container/heap"
	"fmt"
)

// An AStar object allows shortest path searches by variants of the A*
// algorithm.
//
// The variants are determined by the specific search method and heuristic used.
//
// Construct with NewAStar
type AStar struct {
	Graph  WeightedAdjacencyList // input graph
	Result *WeightedFromTree
	// heap values
	r []rNode
	// test instrumentation
	ndVis, arcVis int
}

// NewAStar creates an AStar struct that allows shortest path searches.
//
// Argument g is the graph to be searched, as a weighted adjacency list.
// As usual for AStar, arc weights must be non-negative.
// Graphs may be directed or undirected.
//
// The graph g will not be modified by any AStar methods. NewAStar initializes
// the AStar object for the length (number of nodes) of g. If you add nodes
// to your graph, abandon any previously created Dijkstra object and call
// NewAStar again.
//
// Searches on a single AStar object can be run consecutively but not
// concurrently. Searches can be run concurrently however, on AStar objects
// obtained with separate calls to NewAStar, even with the same graph argument
// to NewAStar.
func NewAStar(g WeightedAdjacencyList) *AStar {
	r := make([]rNode, len(g))
	for i := range r {
		r[i].nx = i
	}
	return &AStar{
		Graph:  g,
		Result: NewWeightedFromTree(len(g)),
		r:      r,
	}
}

// rNode holds data for a "reached" node
type rNode struct {
	nx    int
	state int8 // state constants defined below
	//	prevNode *rNode  // chain encodes path back to start
	//	prevArc  float64 // Arc weight from prevNode to the node of this struct
	//	g        float64 // "g" best known path distance from start node
	f float64 // "g+h", path dist + heuristic estimate
	//	n        int     // number of nodes in path
	fx int // heap.Fix index
}

// for rNode.state
const (
	unreached = 0
	reached   = 1
	open      = 1
	closed    = 2
)

type openHeap []*rNode

// A Heuristic is defined on a specific end node.  The function
// returns an estimate of the path distance from node argument
// "from" to the end node.  Two subclasses of heuristics are "admissable"
// and "monotonic."
//
// Admissable means the value returned is guaranteed to be less than or
// equal to the actual shortest path distance from the node to end.
//
// An admissable estimate may further be monotonic.
// Monotonic means that for any neighboring nodes A and B with half arc aB
// leading from A to B, and for heuristic h defined on some end node, then
// h(A) <= aB.ArcWeight + h(B).
//
// See AStarA for additional notes on implementing heuristic functions for
// AStar search methods.
type Heuristic func(from int) float64

// Admissable returns true if heuristic h is admissable on graph g relative to
// the given end node.
//
// If h is inadmissable, the string result describes a counter example.
func (h Heuristic) Admissable(g WeightedAdjacencyList, end int) (bool, string) {
	// invert graph
	inv := make(WeightedAdjacencyList, len(g))
	for from, nbs := range g {
		for _, nb := range nbs {
			inv[nb.To] = append(inv[nb.To], Half{from, nb.ArcWeight})
		}
	}
	// run dijkstra
	d := NewDijkstra(inv)
	// Dijkstra.AllPaths takes a start node but after inverting the graph
	// argument end now represents the start node of the inverted graph.
	d.AllPaths(end)
	// compare h to found shortest paths
	for n := range inv {
		if !(h(n) <= d.Result.Paths[n].Dist) {
			return false, fmt.Sprintf("h(%d) not <= found shortest path"+
				":  %g not <= %g", n, h(n), d.Result.Paths[n].Dist)
		}
	}
	return true, ""
}

// Monotonic returns true if heuristic h is monotonic on weighted graph g.
//
// If h is non-monotonic, the string result describes a counter example.
func (h Heuristic) Monotonic(g WeightedAdjacencyList) (bool, string) {
	// precompute
	hv := make([]float64, len(g))
	for n := range g {
		hv[n] = h(n)
	}
	// iterate over all edges
	for from, nbs := range g {
		for _, nb := range nbs {
			if !(hv[from] <= nb.ArcWeight+hv[nb.To]) {
				return false, fmt.Sprintf("h(%d) not <= arc weight + h(%d)"+
					":  %g not <= %g + %g",
					from, nb.To, hv[from], nb.ArcWeight, hv[nb.To])
			}
		}
	}
	return true, ""
}

// AStarA finds a path between two nodes.
//
// AStarA implements both algorithm A and algorithm A*.  The difference in the
// two algorithms is strictly in the heuristic estimate returned by argument h.
// If h is an "admissable" heuristic estimate, then the algorithm is termed A*,
// otherwise it is algorithm A.
//
// Like Dijkstra's algorithm, AStarA with an admissable heuristic finds the
// shortest path between start and end.  AStarA generally runs faster than
// Dijkstra though, by using the heuristic distance estimate.
//
// AStarA with an inadmissable heuristic becomes algorithm A.  Algorithm A
// will find a path, but it is not guaranteed to be the shortest path.
// The heuristic still guides the search however, so a nearly admissable
// heuristic is likely to find a very good path, if not the best.  Quality
// of the path returned degrades gracefully with the quality of the heuristic.
//
// The heuristic function h should ideally be fairly inexpensive.  AStarA
// may call it more than once for the same node, especially as graph density
// increases.  In some cases it may be worth the effort to memoize or
// precompute values.
//
// If AStarA finds a path it returns true.  The path can be decoded from
// a.Result.
func (a *AStar) AStarA(start, end int, h Heuristic) bool {
	// NOTE: AStarM is largely duplicate code.

	// reset from any previous run
	a.ndVis = 0
	a.arcVis = 0
	a.Result.reset()
	for i := range a.r {
		a.r[i].state = unreached
	}

	// start node is reached initially
	cr := &a.r[start]
	cr.state = reached
	cr.f = h(start) // total path estimate is estimate from start
	rp := a.Result.Paths
	rp[start].Len = 1  // path length at start is 1 node
	rp[start].Dist = 0 // distance at start is 0
	// oh is a heap of nodes "open" for exploration.  nodes go on the heap
	// when they get an initial or new "g" path distance, and therefore a
	// new "f" which serves as priority for exploration.
	oh := openHeap{cr}
	for len(oh) > 0 {
		bestPath := heap.Pop(&oh).(*rNode)
		bestNode := bestPath.nx
		if bestNode == end {
			return true
		}
		bp := &rp[bestNode]
		nextLen := bp.Len + 1
		for _, nb := range a.Graph[bestNode] {
			alt := &a.r[nb.To]
			ap := &rp[alt.nx]
			g := bp.Dist + nb.ArcWeight // "g" path distance from start
			if alt.state == reached {
				if g > ap.Dist {
					// candidate path to nb is longer than some alternate path
					continue
				}
				if g == ap.Dist && nextLen >= ap.Len {
					// candidate path has identical length of some alternate
					// path but it takes no fewer hops.
					continue
				}
				// cool, we found a better way to get to this node.
				// record new path data for this node and
				// update alt with new data and make sure it's on the heap.
				*ap = WeightedPathEnd{
					From: FromHalf{bestNode, nb.ArcWeight},
					Dist: g,
					Len:  nextLen,
				}
				alt.f = g + h(nb.To)
				if alt.fx < 0 {
					heap.Push(&oh, alt)
				} else {
					heap.Fix(&oh, alt.fx)
				}
			} else {
				// bestNode being reached for the first time.
				*ap = WeightedPathEnd{
					From: FromHalf{bestNode, nb.ArcWeight},
					Dist: g,
					Len:  nextLen,
				}
				alt.f = g + h(nb.To)
				alt.state = reached
				heap.Push(&oh, alt) // and it's now open for exploration
			}
		}
	}
	return false // no path
}

// AStarAPath finds a single shortest path.
//
// Returned is the path and distance as returned by WeightedFromTree.PathTo.
func (a *AStar) AStarAPath(start, end int, h Heuristic) ([]Half, float64) {
	a.AStarA(start, end, h)
	return a.Result.PathTo(end)
}

// AStarM is AStarA optimized for monotonic heuristic estimates.
//
// See AStarA for general usage.  See Heuristic for notes on monotonicity.
func (a *AStar) AStarM(start, end int, h Heuristic) bool {
	// NOTE: AStarM is largely code duplicated from AStarA.
	// Differences are noted in comments in this method.

	a.ndVis = 0
	a.arcVis = 0
	a.Result.reset()
	for i := range a.r {
		a.r[i].state = unreached
	}
	cr := &a.r[start]

	// difference from AStarA:
	// instead of a bit to mark a reached node, there are two states,
	// open and closed. open marks nodes "open" for exploration.
	// nodes are marked open as they are reached, then marked
	// closed as they are found to be on the best path.
	cr.state = open

	cr.f = h(start)
	rp := a.Result.Paths
	rp[start].Len = 1
	rp[start].Dist = 0
	oh := openHeap{cr}
	for len(oh) > 0 {
		bestPath := heap.Pop(&oh).(*rNode)
		bestNode := bestPath.nx
		if bestNode == end {
			return true
		}

		// difference from AStarA:
		// move nodes to closed list as they are found to be best so far.
		bestPath.state = closed

		bp := &rp[bestNode]
		nextLen := bp.Len + 1
		for _, nb := range a.Graph[bestNode] {
			alt := &a.r[nb.To]

			// difference from AStarA:
			// Monotonicity means that f cannot be improved.
			if alt.state == closed {
				continue
			}

			ap := &rp[alt.nx]
			g := bp.Dist + nb.ArcWeight

			// difference from AStarA:
			// test for open state, not just reached
			if alt.state == open {

				if g > ap.Dist {
					continue
				}
				if g == ap.Dist && nextLen >= ap.Len {
					continue
				}
				*ap = WeightedPathEnd{
					From: FromHalf{bestNode, nb.ArcWeight},
					Dist: g,
					Len:  nextLen,
				}
				alt.f = g + h(nb.To)

				// difference from AStarA:
				// we know alt was on the heap because we found it marked open
				heap.Fix(&oh, alt.fx)
			} else {
				*ap = WeightedPathEnd{
					From: FromHalf{bestNode, nb.ArcWeight},
					Dist: g,
					Len:  nextLen,
				}
				alt.f = g + h(nb.To)

				// difference from AStarA:
				// nodes are opened when first reached
				alt.state = open
				heap.Push(&oh, alt)
			}
		}
	}
	return false
}

// AStarMPath finds a single shortest path.
//
// Returned is the path and distance as returned by WeightedFromTree.PathTo.
func (a *AStar) AStarMPath(start, end int, h Heuristic) ([]Half, float64) {
	a.AStarM(start, end, h)
	return a.Result.PathTo(end)
}

// implement container/heap
func (h openHeap) Len() int           { return len(h) }
func (h openHeap) Less(i, j int) bool { return h[i].f < h[j].f }
func (h openHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].fx = i
	h[j].fx = j
}
func (p *openHeap) Push(x interface{}) {
	h := *p
	fx := len(h)
	h = append(h, x.(*rNode))
	h[fx].fx = fx
	*p = h
}

func (p *openHeap) Pop() interface{} {
	h := *p
	last := len(h) - 1
	*p = h[:last]
	h[last].fx = -1
	return h[last]
}
