// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"container/heap"
	"fmt"
	"math"
)

// An AStar object allows shortest path searches by variants of the A*
// algorithm.
//
// The variants are determined by the specific search method and heuristic used.
type AStar struct {
	g WeightedAdjacencyList // input graph
	// r is a list of all nodes reached so far.
	// the chain of nodes following the prev member represents the
	// best path found so far from the start to this node.
	r []rNode
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
		r[i].nd = i
	}
	return &AStar{g: g, r: r}
}

// rNode holds data for a "reached" node
type rNode struct {
	nd       int
	state    int8    // state constants defined below
	prevNode *rNode  // chain encodes path back to start
	prevArc  float64 // Arc weight from prevNode to the node of this struct
	g        float64 // "g" best known path distance from start node
	f        float64 // "g+h", path dist + heuristic estimate
	n        int     // number of nodes in path
	fx       int     // heap.Fix index
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
// The slice result represents the found path with a sequence of half arcs.
// If no path exists from start to end the slice result will be nil. For the
// first element, representing the start node, the arc weight is meaningless
// and will be WeightedFromTree.NoPath. The total path distance is also
// returned.  Path distance is the sum of arc weights, excluding the
// meaningless arc weight of the first Half.
//
// The heuristic function h should ideally be fairly inexpensive.  AStarA
// may call it more than once for the same node, especially as graph density
// increases.  In some cases it may be worth the effort to memoize values.
// Faster yet would be a fully precomputed lookup table, but typically h is
// needed for a rather small fraction of the nodes in the graph.  Construction
// of a complete lookup table will often not be worthwhile.  Profile to see
// if it is even important; benchmark different approaches to find the best.
func (a *AStar) AStarA(start, end int, h Heuristic) ([]Half, float64) {
	// NOTE: AStarM is largely duplicate code.

	// start node is reached initially
	p := &a.r[start]
	p.state = reached
	p.f = h(start) // total path estimate is estimate from start
	p.n = 1        // path length is 1 node
	// oh is a heap of nodes "open" for exploration.  nodes go on the heap
	// when they get an initial or new "g" path distance, and therefore a
	// new "f" which serves as priority for exploration.
	oh := openHeap{p}
	for len(oh) > 0 {
		bestPath := heap.Pop(&oh).(*rNode)
		bestNode := bestPath.nd
		if bestNode == end {
			// done
			dist := bestPath.g
			i := bestPath.n
			path := make([]Half, i)
			for i > 0 {
				i--
				path[i] = Half{To: bestPath.nd, ArcWeight: bestPath.prevArc}
				bestPath = bestPath.prevNode
			}
			return path, dist
		}
		for _, nb := range a.g[bestNode] {
			ed := nb.ArcWeight
			nd := nb.To
			g := bestPath.g + ed
			if alt := &a.r[nd]; alt.state == reached {
				if g > alt.g {
					// new path to nd is longer than some alternate path
					continue
				}
				if g == alt.g && bestPath.n+1 >= alt.n {
					// new path has identical length of some alternate path
					// but it takes more hops.  stick with fewer nodes in path.
					continue
				}
				// cool, we found a better way to get to this node.
				// update alt with new data and make sure it's on the heap.
				alt.prevNode = bestPath
				alt.prevArc = ed
				alt.g = g
				alt.f = g + h(nd)
				alt.n = bestPath.n + 1
				if alt.fx < 0 {
					heap.Push(&oh, alt)
				} else {
					heap.Fix(&oh, alt.fx)
				}
			} else {
				// bestNode being reached for the first time.
				alt.state = reached
				alt.prevNode = bestPath
				alt.prevArc = ed
				alt.g = g
				alt.f = g + h(nd)
				alt.n = bestPath.n + 1
				heap.Push(&oh, alt) // and it's now open for exploration
			}
		}
	}
	return nil, math.Inf(1) // no path
}

// AStarM is AStarA optimized for monotonic heuristic estimates.
//
// See AStarA for general usage.  See Heuristic for notes on monotonicity.
func (a *AStar) AStarM(start, end int, h Heuristic) ([]Half, float64) {
	p := &a.r[start]
	p.f = h(start) // total path estimate is estimate from start
	p.n = 1        // path length is 1 node

	// NOTE: AStarM is largely code duplicated from AStarA.
	// Differences are noted in comments in this method.

	// difference from AStarA:
	// instead of a bit to mark a reached node, there are two states,
	// open and closed. open marks nodes "open" for exploration.
	// nodes are marked open as they are reached, then marked
	// closed as they are found to be on the best path.
	p.state = open

	oh := openHeap{p}
	for len(oh) > 0 {
		bestPath := heap.Pop(&oh).(*rNode)
		bestNode := bestPath.nd
		if bestNode == end {
			// done
			dist := bestPath.g
			i := bestPath.n
			path := make([]Half, i)
			for i > 0 {
				i--
				path[i] = Half{To: bestPath.nd, ArcWeight: bestPath.prevArc}
				bestPath = bestPath.prevNode
			}
			return path, dist
		}

		// difference from AStarA:
		// move nodes to closed list as they are found to be best so far.
		bestPath.state = closed

		for _, nb := range a.g[bestNode] {
			ed := nb.ArcWeight
			nd := nb.To

			// difference from AStarA:
			// Monotonicity means that f cannot be improved.
			if a.r[nd].state == closed {
				continue
			}

			g := bestPath.g + ed
			if alt := &a.r[nd]; alt.state == open {
				if g > alt.g {
					// new path to nd is longer than some alternate path
					continue
				}
				if g == alt.g && bestPath.n+1 >= alt.n {
					// new path has identical length of some alternate path
					// but it takes more hops.  stick with fewer nodes in path.
					continue
				}
				// cool, we found a better way to get to this node.
				// update alt with new data and make sure it's on the heap.
				alt.prevNode = bestPath
				alt.prevArc = ed
				alt.g = g
				alt.f = g + h(nd)
				alt.n = bestPath.n + 1

				// difference from AStarA:
				// we know alt was on the heap because we found it marked open
				heap.Fix(&oh, alt.fx)
			} else {
				// bestNode being reached for the first time.
				alt.state = open
				alt.prevNode = bestPath
				alt.prevArc = ed
				alt.g = g
				alt.f = g + h(nd)
				alt.n = bestPath.n + 1
				heap.Push(&oh, alt) // and it's now open for exploration
			}
		}
	}
	return nil, math.Inf(1) // no path
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
