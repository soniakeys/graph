// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
)

// An AStar object allows shortest path searches by variants of the A*
// algorithm.
//
// The variants are determined by the specific search method and heuristic used.
//
// Construct with NewAStar
type AStar struct {
	Graph  LabeledAdjacencyList // input graph
	Weight WeightFunc           // input weight function
	Tree   FromList             // result paths
	Dist   []float64            // distances for result paths
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
// the AStar object for the order (number of nodes) of g. If you add nodes
// to your graph, abandon any previously created Dijkstra object and call
// NewAStar again.
//
// Searches on a single AStar object can be run consecutively but not
// concurrently. Searches can be run concurrently however, on AStar objects
// obtained with separate calls to NewAStar, even with the same graph argument
// to NewAStar.
func NewAStar(g LabeledAdjacencyList, w WeightFunc) *AStar {
	r := make([]rNode, len(g))
	for i := range r {
		r[i].nx = NI(i)
	}
	return &AStar{
		Graph:  g,
		Weight: w,
		Tree:   NewFromList(len(g)),
		Dist:   make([]float64, len(g)),
		r:      r,
	}
}

// Reset zeros results from any previous search.
//
// It leaves the graph and weight function initialized and otherwise prepares
// the receiver for another search.
func (a *AStar) Reset() {
	a.ndVis = 0
	a.arcVis = 0
	a.Tree.reset()
	for i := range a.r {
		a.r[i].state = unreached
	}
	for i := range a.Dist {
		a.Dist[i] = 0
	}
}

// AllPaths finds shortest paths from start to all nodes reachable
// from start.
//
// AllPaths returns number of paths found, equivalent to the number of nodes
// reached, including the path ending at start.  Path results are left in
// d.Result.
// AllPaths finds shortest paths from start to all nodes reachable
// from start.
//
// AllPaths returns number of paths found, equivalent to the number of nodes
// reached, including the path ending at start.  Path results are left in
// d.Result.
// rNode holds data for a "reached" node
type rNode struct {
	nx    NI
	state int8    // state constants defined below
	f     float64 // "g+h", path dist + heuristic estimate
	fx    int     // heap.Fix index
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
type Heuristic func(from NI) float64

// Admissable returns true if heuristic h is admissable on graph g relative to
// the given end node.
//
// If h is inadmissable, the string result describes a counter example.
func (h Heuristic) Admissable(g LabeledAdjacencyList, w WeightFunc, end NI) (bool, string) {
	// invert graph
	inv := make(LabeledAdjacencyList, len(g))
	for from, nbs := range g {
		for _, nb := range nbs {
			inv[nb.To] = append(inv[nb.To],
				Half{To: NI(from), Label: nb.Label})
		}
	}
	// run dijkstra
	d := NewDijkstra(inv, w)
	// Dijkstra.AllPaths takes a start node but after inverting the graph
	// argument end now represents the start node of the inverted graph.
	d.AllPaths(end)
	// compare h to found shortest paths
	for n := range inv {
		if d.Tree.Paths[n].Len == 0 {
			continue // no path, any heuristic estimate is fine.
		}
		if !(h(NI(n)) <= d.Dist[n]) {
			return false, fmt.Sprintf("h(%d) = %g, "+
				"required to be <= found shortest path (%g)",
				n, h(NI(n)), d.Dist[n])
		}
	}
	return true, ""
}

// Monotonic returns true if heuristic h is monotonic on weighted graph g.
//
// If h is non-monotonic, the string result describes a counter example.
func (h Heuristic) Monotonic(g LabeledAdjacencyList, w WeightFunc) (bool, string) {
	// precompute
	hv := make([]float64, len(g))
	for n := range g {
		hv[n] = h(NI(n))
	}
	// iterate over all edges
	for from, nbs := range g {
		for _, nb := range nbs {
			arcWeight := w(nb.Label)
			if !(hv[from] <= arcWeight+hv[nb.To]) {
				return false, fmt.Sprintf("h(%d) = %g, "+
					"required to be <= arc weight + h(%d) (= %g + %g = %g)",
					from, hv[from],
					nb.To, arcWeight, hv[nb.To], arcWeight+hv[nb.To])
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
func (a *AStar) AStarA(start, end NI, h Heuristic) bool {
	// NOTE: AStarM is largely duplicate code.

	// start node is reached initially
	cr := &a.r[start]
	cr.state = reached
	cr.f = h(start) // total path estimate is estimate from start
	rp := a.Tree.Paths
	rp[start] = PathEnd{Len: 1, From: -1} // path length at start is 1 node
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
			// "g" path distance from start
			g := a.Dist[bestNode] + a.Weight(nb.Label)
			if alt.state == reached {
				if g > a.Dist[nb.To] {
					// candidate path to nb is longer than some alternate path
					continue
				}
				if g == a.Dist[nb.To] && nextLen >= ap.Len {
					// candidate path has identical length of some alternate
					// path but it takes no fewer hops.
					continue
				}
				// cool, we found a better way to get to this node.
				// record new path data for this node and
				// update alt with new data and make sure it's on the heap.
				*ap = PathEnd{From: bestNode, Len: nextLen}
				a.Dist[nb.To] = g
				alt.f = g + h(nb.To)
				if alt.fx < 0 {
					heap.Push(&oh, alt)
				} else {
					heap.Fix(&oh, alt.fx)
				}
			} else {
				// bestNode being reached for the first time.
				*ap = PathEnd{From: bestNode, Len: nextLen}
				a.Dist[nb.To] = g
				alt.f = g + h(nb.To)
				alt.state = reached
				heap.Push(&oh, alt) // and it's now open for exploration
			}
		}
	}
	return false // no path
}

// AStarAPath finds a single shortest path using the AStarA algorithm.
//
// See documentation on the AStarA method of the AStar type.
//
// Returned is the path and distance as returned by FromList.PathTo.
func (g LabeledAdjacencyList) AStarAPath(start, end NI, h Heuristic, w WeightFunc) ([]NI, float64) {
	a := NewAStar(g, w)
	a.AStarA(start, end, h)
	return a.Tree.PathTo(end, nil), a.Dist[end]
}

// AStarM is AStarA optimized for monotonic heuristic estimates.
//
// Note that this function requires a monotonic heuristic.  Results will
// not be meaningful if argument h is non-monotonic.
//
// See AStarA for general usage.  See Heuristic for notes on monotonicity.
func (a *AStar) AStarM(start, end NI, h Heuristic) bool {
	// NOTE: AStarM is largely code duplicated from AStarA.
	// Differences are noted in comments in this method.

	cr := &a.r[start]

	// difference from AStarA:
	// instead of a bit to mark a reached node, there are two states,
	// open and closed. open marks nodes "open" for exploration.
	// nodes are marked open as they are reached, then marked
	// closed as they are found to be on the best path.
	cr.state = open

	cr.f = h(start)
	rp := a.Tree.Paths
	rp[start] = PathEnd{Len: 1, From: -1}
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
			g := a.Dist[bestNode] + a.Weight(nb.Label)

			// difference from AStarA:
			// test for open state, not just reached
			if alt.state == open {

				if g > a.Dist[nb.To] {
					continue
				}
				if g == a.Dist[nb.To] && nextLen >= ap.Len {
					continue
				}
				*ap = PathEnd{From: bestNode, Len: nextLen}
				a.Dist[nb.To] = g
				alt.f = g + h(nb.To)

				// difference from AStarA:
				// we know alt was on the heap because we found it marked open
				heap.Fix(&oh, alt.fx)
			} else {
				*ap = PathEnd{From: bestNode, Len: nextLen}
				a.Dist[nb.To] = g
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

// AStarMPath finds a single shortest path using the AStarM algorithm.
//
// See documentation on the AStarM method of the AStar type.
//
// Returned is the path and distance as returned by FromList.PathTo.
func (g LabeledAdjacencyList) AStarMPath(start, end NI, h Heuristic, w WeightFunc) ([]NI, float64) {
	a := NewAStar(g, w)
	a.AStarM(start, end, h)
	return a.Tree.PathTo(end, nil), a.Dist[end]
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

// A BellmanFord object allows shortest path searches using the
// Bellman-Ford-Moore algorithm.
type BellmanFord struct {
	Graph  LabeledAdjacencyList
	Weight WeightFunc
	Tree   FromList
	Dist   []float64
}

// NewBellmanFord creates a BellmanFord object that allows shortest path
// searches using the Bellman-Ford-Moore algorithm.
//
// Argument g is the graph to be searched, as a labeled adjacency list.
// WeightFunc w must translate arc labels to arc weights.
// Negative arc weights are allowed as long as there are no negative cycles.
// Graphs may be directed or undirected.  Loops and parallel arcs are
// allowed.
//
// The graph g will not be modified by any BellmanFord methods.  NewBellmanFord
// initializes the BellmanFord object for the order (number of nodes) of g.
// If you add nodes to your graph, abandon any previously created BellmanFord
// object and call NewBellmanFord again.
//
// Searches on a single BellmanFord object can be run consecutively but not
// concurrently.  Searches can be run concurrently however, on BellmanFord
// objects obtained with separate calls to NewBellmanFord, even with the same
// graph argument to NewBellmanFord.
func NewBellmanFord(g LabeledAdjacencyList, w WeightFunc) *BellmanFord {
	return &BellmanFord{
		Graph:  g,
		Weight: w,
		Tree:   NewFromList(len(g)),
		Dist:   make([]float64, len(g)),
	}
}

// Reset zeros results from any previous search.
//
// It leaves the graph and weight function initialized and otherwise prepares
// the receiver for another search.
func (b *BellmanFord) Reset() {
	b.Tree.reset()
	for i := range b.Dist {
		b.Dist[i] = 0
	}
}

// Start runs the BellmanFord algorithm to finds shortest paths from start
// to all nodes reachable from start.
//
// The algorithm allows negative edge weights but not negative cycles.
// Start returns true if the algorithm completes successfully.  In this case
// FromList b.Tree will encode the shortest paths found.
//
// Start returns false in the case that it encounters a negative cycle
// reachable from start.  In this case values in b.Result are meaningless.
//
// Negative cycles are only detected when reachable from start.  A negative
// cycle not reachable from start will not prevent the algorithm from finding
// shortest paths reachable from start.
func (b *BellmanFord) Start(start NI) (ok bool) {
	inf := math.Inf(1)
	for i := range b.Dist {
		b.Dist[i] = inf
	}
	rp := b.Tree.Paths
	rp[start] = PathEnd{Len: 1, From: -1}
	b.Dist[start] = 0
	for _ = range b.Graph[1:] {
		imp := false
		for from, nbs := range b.Graph {
			fp := &rp[from]
			d1 := b.Dist[from]
			for _, nb := range nbs {
				d2 := d1 + b.Weight(nb.Label)
				to := &rp[nb.To]
				// TODO improve to break ties
				if fp.Len > 0 && d2 < b.Dist[nb.To] {
					*to = PathEnd{From: NI(from), Len: fp.Len + 1}
					b.Dist[nb.To] = d2
					imp = true
				}
			}
		}
		if !imp {
			break
		}
	}
	for from, nbs := range b.Graph {
		d1 := b.Dist[from]
		for _, nb := range nbs {
			if d1+b.Weight(nb.Label) < b.Dist[nb.To] {
				return false // negative cycle
			}
		}
	}
	return true
}

// NegativeCycle returns true if the graph contains any negative cycle.
//
// Path information is not computed.
//
// Note the sense of the returned value is opposite that of BellmanFord.Start().
func (b *BellmanFord) NegativeCycle() bool {
	for i := range b.Dist {
		b.Dist[i] = 0
	}
	for _ = range b.Graph[1:] {
		imp := false
		for from, nbs := range b.Graph {
			d1 := b.Dist[from]
			for _, nb := range nbs {
				d2 := d1 + b.Weight(nb.Label)
				if d2 < b.Dist[nb.To] {
					b.Dist[nb.To] = d2
					imp = true
				}
			}
		}
		if !imp {
			break
		}
	}
	for from, nbs := range b.Graph {
		d1 := b.Dist[from]
		for _, nb := range nbs {
			if d1+b.Weight(nb.Label) < b.Dist[nb.To] {
				return true // negative cycle
			}
		}
	}
	return false
}

// BreadthFirst associates a graph with a result object for returning
// results from breadth first searches and traversals.
//
// Construct with NewBreadthFirst.
//
// With all methods that traverse the graph, if the Rand member is left nil,
// nodes are visited in the order they exist as to-arcs in the graph.
// If Rand is assigned a random number generator, nodes are visited in random
// order at each level.  See example at Traverse method.
//
// The search methods set Result.Paths and Result.MaxLen but not Result.Leaves.
type BreadthFirst struct {
	Graph  AdjacencyList
	Rand   *rand.Rand
	Result FromList
}

// NewBreadthFirst creates a BreadthFirst object.
//
// Argument g is the graph to be searched, as an adjacency list.
// Graphs may be directed or undirected.
//
// The graph g will not be modified by any BreadthFirst methods.
// NewBreadthFirst initializes the BreadthFirst object for the order
// (number of nodes) of g.  If you add nodes to your graph, abandon any
// previously created BreadthFirst object and call NewBreadthFirst again.
//
// Searches on a single BreadthFirst object can be run consecutively but not
// concurrently.  Searches can be run concurrently however, on BreadthFirst
// objects obtained with separate calls to NewBreadthFirst, even with the same
// graph argument to NewBreadthFirst.
func NewBreadthFirst(g AdjacencyList) *BreadthFirst {
	return &BreadthFirst{
		Graph:  g,
		Result: NewFromList(len(g)),
	}
}

// BreadthFirstPath finds a single path from start to end with a minimum
// number of nodes.
//
// Returned is the path as list of nodes.
// The result is nil if no path was found.
func (g AdjacencyList) BreadthFirstPath(start, end NI) []NI {
	b := NewBreadthFirst(g)
	b.Traverse(start, func(n NI) bool { return n != end })
	return b.Result.PathTo(end, nil)
}

// Path finds a single path from start to end with a minimum number of nodes.
//
// Path returns true if a path exists, false if not.  The path can be recovered
// from b.Result.
func (b *BreadthFirst) Path(start, end NI) bool {
	b.Traverse(start, func(n NI) bool { return n != end })
	return b.Result.Paths[end].Len > 0
}

// AllPaths finds paths from start to all nodes reachable from start that
// have a minimum number of nodes.
//
// AllPaths returns number of paths found, equivalent to the number of nodes
// reached, including the path ending at start.  Path results are left in
// d.Result.
func (b *BreadthFirst) AllPaths(start NI) int {
	return b.Traverse(start, func(NI) bool { return true })
}

// A Visitor function is an argument to graph traversal methods.
//
// Graph traversal methods call the visitor function for each node visited.
// The argument n is the node being visited.  If the visitor function
// returns true, the traversal will continue.  If the visitor function
// returns false, the traversal will terminate immediately.
type Visitor func(n NI) (ok bool)

// Traverse traverses a graph in breadth first order starting from node start.
//
// Traverse calls the visitor function v for each node.  If v returns true,
// the traversal will continue.  If v returns false, the traversal will
// terminate immediately.  Traverse updates b.Result.Paths before calling v.
//
// Traverse returns the number of nodes successfully visited; if search is
// is terminated by a false return from v, that node is not counted.
func (b *BreadthFirst) Traverse(start NI, v Visitor) int {
	b.Result.reset()
	rp := b.Result.Paths
	nVisitOk := 0 // accumulated for a return value
	// the frontier consists of nodes all at the same level
	frontier := []NI{start}
	level := 1
	// assign path when node is put on frontier,
	rp[start] = PathEnd{Len: level, From: -1}
	for {
		b.Result.MaxLen = level
		level++
		var next []NI
		if b.Rand == nil {
			for _, n := range frontier {
				if !v(n) { // visit nodes as they come off frontier
					return -1
				}
				nVisitOk++
				for _, nb := range b.Graph[n] {
					if rp[nb].Len == 0 {
						next = append(next, nb)
						rp[nb] = PathEnd{From: n, Len: level}
					}
				}
			}
		} else { // take nodes off frontier at random
			for _, i := range b.Rand.Perm(len(frontier)) {
				n := frontier[i]
				// remainder of block same as above
				if !v(n) {
					return -1
				}
				nVisitOk++
				for _, nb := range b.Graph[n] {
					if rp[nb].Len == 0 {
						next = append(next, nb)
						rp[nb] = PathEnd{From: n, Len: level}
					}
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
	}
	return nVisitOk
}

// BreadthFirst2 methods implement a direction-optimized strategy.
// Code is experimental and currently is no faster than basic BreadthFirst.
type BreadthFirst2 struct {
	To, From AdjacencyList
	MArc     int
	Result   FromList
}

func NewBreadthFirst2(to, from AdjacencyList, ma int) *BreadthFirst2 {
	return &BreadthFirst2{
		To:     to,
		From:   from,
		MArc:   ma,
		Result: NewFromList(len(to)),
	}
}

// BreadthFirst2Path finds a single path from start to end with a minimum
// number of nodes using a direction optimizing algorithm.
//
// It is experimental and not currently recommended over BreadthFirstPath
// as it is currently no faster and offers no benefits.
//
// Returned is the path as list of nodes.
// The result is nil if no path was found.
func (g AdjacencyList) BreadthFirst2Path(start, end NI) []NI {
	t, ma := Directed{g}.Transpose()
	b := NewBreadthFirst2(g, t.AdjacencyList, ma)
	b.Traverse(start, func(n NI) bool { return n != end })
	return b.Result.PathTo(end, nil)
}

func (b *BreadthFirst2) Path(start, end NI) bool {
	b.Traverse(start, func(n NI) bool { return n != end })
	return b.Result.Paths[end].Len > 0
}

func (b *BreadthFirst2) AllPaths(start NI) int {
	return b.Traverse(start, func(NI) bool { return true })
}

func (b *BreadthFirst2) Traverse(start NI, v Visitor) int {
	b.Result.reset()
	rp := b.Result.Paths
	level := 1
	rp[start] = PathEnd{Len: level, From: -1}
	if !v(start) {
		b.Result.MaxLen = level
		return -1
	}
	nReached := 1 // accumulated for a return value
	// the frontier consists of nodes all at the same level
	frontier := []NI{start}
	mf := len(b.To[start])         // number of arcs leading out from frontier
	ctb := b.MArc / 10             // threshold change from top-down to bottom-up
	k14 := 14 * b.MArc / len(b.To) // 14 * mean degree
	cbt := len(b.To) / k14         // threshold change from bottom-up to top-down
	//	var fBits, nextb big.Int
	fBits := make([]bool, len(b.To))
	nextb := make([]bool, len(b.To))
	zBits := make([]bool, len(b.To))
	for {
		// top down step
		level++
		var next []NI
		for _, n := range frontier {
			for _, nb := range b.To[n] {
				if rp[nb].Len == 0 {
					rp[nb] = PathEnd{From: n, Len: level}
					if !v(nb) {
						b.Result.MaxLen = level
						return -1
					}
					next = append(next, nb)
					nReached++
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
		if mf > ctb {
			// switch to bottom up!
		} else {
			// stick with top down
			continue
		}
		// convert frontier representation
		nf := 0 // number of vertices on the frontier
		for _, n := range frontier {
			//			fBits.SetBit(&fBits, n, 1)
			fBits[n] = true
			nf++
		}
	bottomUpLoop:
		level++
		nNext := 0
		for n := range b.From {
			if rp[n].Len == 0 {
				for _, nb := range b.From[n] {
					//					if fBits.Bit(nb) == 1 {
					if fBits[nb] {
						rp[n] = PathEnd{From: nb, Len: level}
						if !v(nb) {
							b.Result.MaxLen = level
							return -1
						}
						//						nextb.SetBit(&nextb, n, 1)
						nextb[n] = true
						nReached++
						nNext++
						break
					}
				}
			}
		}
		if nNext == 0 {
			break
		}
		fBits, nextb = nextb, fBits
		//		nextb.SetInt64(0)
		copy(nextb, zBits)
		nf = nNext
		if nf < cbt {
			// switch back to top down!
		} else {
			// stick with bottom up
			goto bottomUpLoop
		}
		// convert frontier representation
		mf = 0
		frontier = frontier[:0]
		for n := range b.To {
			//			if fBits.Bit(n) == 1 {
			if fBits[n] {
				frontier = append(frontier, NI(n))
				mf += len(b.To[n])
				fBits[n] = false
			}
		}
		//		fBits.SetInt64(0)
	}
	b.Result.MaxLen = level - 1
	return nReached
}

// A DAGPath object allows searches for paths of either shortest or longest
// distance in a directed acyclic graph.
//
// DAGPath methods measure path distance as the sum of arc weights.
// Negative arc weights are allowed.
// Where multiple paths exist with the same distance, the path length
// (number of nodes) is used as a tie breaker.
//
// Construct with NewDAGPath.
type DAGPath struct {
	Graph    LabeledAdjacencyList // input graph
	Ordering []NI                 // topological ordering
	Weight   WeightFunc           // input weight function
	Longest  bool                 // find longest path rather than shortest
	Tree     FromList             // result paths
	Dist     []float64            // in Tree, distances for result paths
}

// NewDAGPath creates a DAGPath object that allows path searches of either
// shortest or longest distance.
//
// Argument g is the graph to be searched, and must be a directed acyclic
// graph.  Argument o must be a topological ordering of g.
//
// The graph g will not be modified by any DAGPath methods.  NewDAGPath
// initializes the DAGPath object for the order (number of nodes) of g.
// If you add nodes to your graph, abandon any previously created DAGPath
// object and call NewDAGPath again.
//
// Searches on a single DAGPath object can be run consecutively but not
// concurrently.  Searches can be run concurrently however, on DAGPath
// objects obtained with separate calls to NewDAGPath, even with the same
// graph argument to NewDAGPath.
func NewDAGPath(g DirectedLabeled, ordering []NI, w WeightFunc, longest bool) *DAGPath {
	return &DAGPath{
		Graph:    g.LabeledAdjacencyList,
		Ordering: ordering,
		Weight:   w,
		Longest:  longest,
		Tree:     NewFromList(len(g.LabeledAdjacencyList)),
		Dist:     make([]float64, len(g.LabeledAdjacencyList)),
	}
}

// Reset zeros results from any previous search.
//
// It leaves members Graph, Ordering, Weight, and Longest initialized and
// otherwise prepares the receiver for another search.
func (d *DAGPath) Reset() {
	d.Tree.reset()
	for i := range d.Dist {
		d.Dist[i] = 0
	}
}

// DAGMinDistPath finds a single shortest path.
//
// Shortest means minimum sum of arc weights.
//
// Returned is the path and distance as returned by FromList.PathTo.
//
// This is a convenience method.  See the DAGPath type for more options.
func (g DirectedLabeled) DAGMinDistPath(start, end NI, w WeightFunc) ([]NI, float64, error) {
	return g.dagPath(start, end, w, false)
}

// DAGMaxDistPath finds a single longest path.
//
// Longest means maximum sum of arc weights.
//
// Returned is the path and distance as returned by FromList.PathTo.
//
// This is a convenience method.  See the DAGPath type for more options.
func (g DirectedLabeled) DAGMaxDistPath(start, end NI, w WeightFunc) ([]NI, float64, error) {
	return g.dagPath(start, end, w, true)
}

func (g DirectedLabeled) dagPath(start, end NI, w WeightFunc, longest bool) ([]NI, float64, error) {
	o, _ := g.Topological()
	if o == nil {
		return nil, 0, fmt.Errorf("not a DAG")
	}
	d := NewDAGPath(g, o, w, longest)
	d.Path(start, end)
	if d.Tree.Paths[end].Len == 0 {
		return nil, 0, fmt.Errorf("no path from %d to %d", start, end)
	}
	return d.Tree.PathTo(end, nil), d.Dist[end], nil
}

// Path finds a single path.
//
// Path returns true if a path exists, false if not. The path can be recovered
// from the receiver.  Path returns as soon as the shortest path to end is
// found; it does not explore the remainder of the graph.
func (d *DAGPath) Path(start, end NI) bool {
	d.search(start, end)
	return d.Tree.Paths[end].Len > 0
}

// AllPaths finds paths from argument start to all nodes reachable from start.
func (d *DAGPath) AllPaths(start NI) (reached int) {
	return d.search(start, -1)
}

func (d *DAGPath) search(start, end NI) (reached int) {
	// search ordering for start
	o := 0
	for d.Ordering[o] != start {
		o++
	}
	var fBetter func(cand, ext float64) bool
	var iBetter func(cand, ext int) bool
	if d.Longest {
		fBetter = func(cand, ext float64) bool { return cand > ext }
		iBetter = func(cand, ext int) bool { return cand > ext }
	} else {
		fBetter = func(cand, ext float64) bool { return cand < ext }
		iBetter = func(cand, ext int) bool { return cand < ext }
	}
	g := d.Graph
	p := d.Tree.Paths
	p[start] = PathEnd{From: -1, Len: 1}
	d.Tree.MaxLen = 1
	leaves := &d.Tree.Leaves
	leaves.SetBit(leaves, int(start), 1)
	reached = 1
	for n := start; n != end; n = d.Ordering[o] {
		if p[n].Len > 0 && len(g[n]) > 0 {
			nDist := d.Dist[n]
			candLen := p[n].Len + 1 // len for any candidate arc followed from n
			for _, to := range g[n] {
				leaves.SetBit(leaves, int(to.To), 1)
				candDist := nDist + d.Weight(to.Label)
				switch {
				case p[to.To].Len == 0: // first path to node to.To
					reached++
				case fBetter(candDist, d.Dist[to.To]): // better distance
				case candDist == d.Dist[to.To] && iBetter(candLen, p[to.To].Len): // same distance but better path length
				default:
					continue
				}
				d.Dist[to.To] = candDist
				p[to.To] = PathEnd{From: n, Len: candLen}
				if candLen > d.Tree.MaxLen {
					d.Tree.MaxLen = candLen
				}
			}
			leaves.SetBit(leaves, int(n), 0)
		}
		o++
		if o == len(d.Ordering) {
			break
		}
	}
	return
}

// A Dijkstra object allows shortest path searches by Dijkstra's algorithm.
//
// Dijkstra methods find paths of shortest distance where distance is the
// sum of arc weights.  Where multiple paths exist with the same distance,
// Dijkstra methods return a path with the minimum number of nodes.
//
// Construct with NewDijkstra.
type Dijkstra struct {
	Graph  LabeledAdjacencyList // input graph
	Weight WeightFunc           // input weight function
	Tree   FromList             // result paths
	Dist   []float64            // in Tree, distances for result paths
	// heap values
	r []tentResult
	// test instrumentation
	ndVis, arcVis int
}

// NewDijkstra creates a Dijkstra object that allows shortest path searches.
//
// Argument g is the graph to be searched, as a weighted adjacency list.
// As usual for Dijkstra's algorithm, arc weights must be non-negative.
// Graphs may be directed or undirected.  Loops and parallel arcs are
// allowed.
//
// The graph g will not be modified by any Dijkstra methods.  NewDijkstra
// initializes the Dijkstra object for the order (number of nodes) of g.
// If you add nodes to your graph, abandon any previously created Dijkstra
// object and call NewDijkstra again.
//
// Searches on a single Dijkstra object can be run consecutively but not
// concurrently.  Searches can be run concurrently however, on Dijkstra
// objects obtained with separate calls to NewDijkstra, even with the same
// graph argument to NewDijkstra.
func NewDijkstra(g LabeledAdjacencyList, w WeightFunc) *Dijkstra {
	r := make([]tentResult, len(g))
	for i := range r {
		r[i].nx = NI(i)
	}
	return &Dijkstra{
		Graph:  g,
		Weight: w,
		Tree:   NewFromList(len(g)),
		Dist:   make([]float64, len(g)),
		r:      r,
	}
}

// Reset zeros results from any previous search.
//
// It leaves the graph and weight function initialized and otherwise prepares
// the receiver for another search.
func (d *Dijkstra) Reset() {
	d.ndVis = 0
	d.arcVis = 0
	d.Tree.reset()
	for i := range d.r {
		d.r[i].done = false
	}
	for i := range d.Dist {
		d.Dist[i] = 0
	}
}

type tentResult struct {
	dist float64 // tentative distance, sum of arc weights
	nx   NI      // slice index, "node id"
	fx   int     // heap.Fix index
	done bool
}

type tent []*tentResult

// DijkstraPath finds a single shortest path.
//
// Returned is the path and distance as returned by FromList.PathTo.
func (g LabeledAdjacencyList) DijkstraPath(start, end NI, w WeightFunc) ([]NI, float64) {
	d := NewDijkstra(g, w)
	d.Path(start, end)
	return d.Tree.PathTo(end, nil), d.Dist[end]
}

// Path finds a single shortest path.
//
// Path returns true if a path exists, false if not. The path can be recovered
// from b.Result.  Path returns as soon as the shortest path to end is found;
// it does not explore the remainder of the graph.
func (d *Dijkstra) Path(start, end NI) bool {
	d.search(start, end)
	return d.Tree.Paths[end].Len > 0
}

// AllPaths finds shortest paths from start to all nodes reachable
// from start.
//
// AllPaths returns number of paths found, equivalent to the number of nodes
// reached, including the path ending at start.  Path results are left in
// d.Result.
func (d *Dijkstra) AllPaths(start NI) (nFound int) {
	return d.search(start, -1)
}

// returns number of nodes reached (= number of shortest paths found)
func (d *Dijkstra) search(start, end NI) (reached int) {
	current := start
	rp := d.Tree.Paths
	rp[current] = PathEnd{Len: 1, From: -1} // path length at start is 1 node
	cr := &d.r[current]
	cr.dist = 0    // distance at start is 0.
	cr.done = true // mark start done.  it skips the heap.
	nDone := 1     // accumulated for a return value
	var t tent
	for current != end {
		nextLen := rp[current].Len + 1
		for _, nb := range d.Graph[current] {
			d.arcVis++
			hr := &d.r[nb.To]
			if hr.done {
				continue // skip nodes already done
			}
			dist := cr.dist + d.Weight(nb.Label)
			vl := rp[nb.To].Len
			visited := vl > 0
			if visited {
				if dist > hr.dist {
					continue // distance is worse
				}
				// tie breaker is a nice touch and doesn't seem to
				// impact performance much.
				if dist == hr.dist && nextLen >= vl {
					continue // distance same, but number of nodes is no better
				}
			}
			// the path through current to this node is shortest so far.
			// record new path data for this node and update tentative set.
			hr.dist = dist
			rp[nb.To].Len = nextLen
			rp[nb.To].From = current
			if visited {
				heap.Fix(&t, hr.fx)
			} else {
				heap.Push(&t, hr)
			}
		}
		d.ndVis++
		if len(t) == 0 {
			return nDone // no more reachable nodes. AllPaths normal return
		}
		// new current is node with smallest tentative distance
		cr = heap.Pop(&t).(*tentResult)
		cr.done = true
		nDone++
		current = cr.nx
		d.Dist[current] = cr.dist // store final distance
	}
	// normal return for single shortest path search
	return -1
}

// tent implements container/heap
func (t tent) Len() int           { return len(t) }
func (t tent) Less(i, j int) bool { return t[i].dist < t[j].dist }
func (t tent) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].fx = i
	t[j].fx = j
}
func (s *tent) Push(x interface{}) {
	nd := x.(*tentResult)
	nd.fx = len(*s)
	*s = append(*s, nd)
}
func (s *tent) Pop() interface{} {
	t := *s
	last := len(t) - 1
	*s = t[:last]
	return t[last]
}
