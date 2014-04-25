// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"container/heap"
	"fmt"
	"math/big"
)

// A BreadthFirst object allows graph traversals and searches in
// breadth first order.
//
// Construct with NewBreadthFirst.
type BreadthFirst struct {
	Graph  AdjacencyList
	Result *FromTree
}

// NewBreadthFirst creates a BreadthFirst object.
//
// Argument g is the graph to be searched, as an adjacency list.
// Graphs may be directed or undirected.
//
// The graph g will not be modified by any BreadthFirst methods.
// NewBreadthFirst initializes the BreadthFirst object for the length
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
		Result: NewFromTree(len(g)),
	}
}

// BreadthFirstPath finds a single path from start to end with a minimum
// number of nodes.
//
// Returned is the path as list of nodes.
// Path returns nil if no path was found.
func BreadthFirstPath(g AdjacencyList, start, end int) []int {
	b := NewBreadthFirst(g)
	b.Traverse(start, func(n int) bool { return n != end })
	return b.Result.PathTo(end)
}

// Path finds a single path from start to end with a minimum number of nodes.
//
// Path returns true if a path exists, false if not.  The path can be recovered
// from b.Result.
func (b *BreadthFirst) Path(start, end int) bool {
	b.Traverse(start, func(n int) bool { return n != end })
	return b.Result.Paths[end].Len > 0
}

// AllPaths finds paths from start to all nodes reachable from start that
// have a minimum number of nodes.
//
// AllPaths returns number of paths found, equivalent to the number of nodes
// reached, including the path ending at start.  Path results are left in
// d.Result.
func (b *BreadthFirst) AllPaths(start int) int {
	return b.Traverse(start, func(int) bool { return true })
}

// A Visitor function is an argument to graph traversal methods.
//
// Graph traversal methods call the visitor function for each node visited.
// The argument n is the node being visited.  If the visitor function
// returns true, the traversal will continue.  If the visitor function
// returns false, the traversal will terminate immediately.
type Visitor func(n int) (ok bool)

// Traverse traverses a graph in breadth first order starting from node start.
//
// Traverse calls the visitor function v for each node.  If v returns true,
// the traversal will continue.  If v returns false, the traversal will
// terminate immediately.  Traverse updates b.Result.Paths before calling v.
// It updates b.Result.MaxLen after calling v, if v returns true.
//
// Traverse returns the number of nodes successfully visited; if search is
// is terminated by a false return from v, that node is not counted.
func (b *BreadthFirst) Traverse(start int, v Visitor) int {
	b.Result.Reset()
	rp := b.Result.Paths
	b.Result.Start = start
	level := 1
	rp[start].Len = level
	if !v(start) {
		b.Result.MaxLen = level
		return -1
	}
	nReached := 1 // accumulated for a return value
	// the frontier consists of nodes all at the same level
	frontier := []int{start}
	for {
		level++
		var next []int
		for _, n := range frontier {
			for _, nb := range b.Graph[n] {
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
	}
	b.Result.MaxLen = level - 1
	return nReached
}

type BreadthFirst2 struct {
	To, From AdjacencyList
	M        int
	Result   *FromTree
}

func NewBreadthFirst2(to, from AdjacencyList, m int) *BreadthFirst2 {
	return &BreadthFirst2{
		To:     to,
		From:   from,
		M:      m,
		Result: NewFromTree(len(to)),
	}
}

func BreadthFirst2Path(g AdjacencyList, start, end int) []int {
	t, m := g.Transpose()
	b := NewBreadthFirst2(g, t, m)
	b.Traverse(start, func(n int) bool { return n != end })
	return b.Result.PathTo(end)
}

func (b *BreadthFirst2) Path(start, end int) bool {
	b.Traverse(start, func(n int) bool { return n != end })
	return b.Result.Paths[end].Len > 0
}

func (b *BreadthFirst2) AllPaths(start int) int {
	return b.Traverse(start, func(int) bool { return true })
}

func (b *BreadthFirst2) Traverse(start int, v Visitor) int {
	b.Result.Reset()
	rp := b.Result.Paths
	b.Result.Start = start
	level := 1
	rp[start].Len = level
	if !v(start) {
		b.Result.MaxLen = level
		return -1
	}
	nReached := 1 // accumulated for a return value
	// the frontier consists of nodes all at the same level
	frontier := []int{start}
	mf := len(b.To[start])      // number of arcs leading out from frontier
	ctb := b.M / 10             // threshold change from top-down to bottom-up
	k14 := 14 * b.M / len(b.To) // 14 * mean degree
	cbt := len(b.To) / k14      // threshold change from bottom-up to top-down
	var fBits, nextb big.Int
	for {
		// top down step
		level++
		var next []int
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
			fBits.SetBit(&fBits, n, 1)
			nf++
		}
	bottomUpLoop:
		level++
		nNext := 0
		for n := range b.From {
			if rp[n].Len == 0 {
				for _, nb := range b.From[n] {
					if fBits.Bit(nb) == 1 {
						rp[n] = PathEnd{From: nb, Len: level}
						if !v(nb) {
							b.Result.MaxLen = level
							return -1
						}
						nextb.SetBit(&nextb, n, 1)
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
		nextb.SetInt64(0)
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
			if fBits.Bit(n) == 1 {
				frontier = append(frontier, n)
				mf += len(b.To[n])
			}
		}
		fBits.SetInt64(0)
	}
	b.Result.MaxLen = level - 1
	return nReached
}

// A Dijkstra object allows shortest path searches by Dijkstra's algorithm.
//
// Dijkstra methods find paths of shortest distance where distance is the
// sum of arc weights.  Where multiple paths exist with the same distance,
// Dijkstra methods return a path with the minimum number of nodes.
//
// Construct with NewDijkstra.
type Dijkstra struct {
	Graph  WeightedAdjacencyList
	Result *WeightedFromTree
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
// initializes the Dijkstra object for the length (number of nodes) of g.
// If you add nodes to your graph, abandon any previously created Dijkstra
// object and call NewDijkstra again.
//
// NewDijkstra calls NewWeightedFromTree.  See documentation there in
// particular for the option to change NoPath before running a search.
//
// Searches on a single Dijkstra object can be run consecutively but not
// concurrently.  Searches can be run concurrently however, on Dijkstra
// objects obtained with separate calls to NewDijkstra, even with the same
// graph argument to NewDijkstra.
func NewDijkstra(g WeightedAdjacencyList) *Dijkstra {
	r := make([]tentResult, len(g))
	for i := range r {
		r[i].nx = i
	}
	return &Dijkstra{
		Graph:  g,
		Result: NewWeightedFromTree(len(g)),
		r:      r,
	}
}

type tentResult struct {
	dist float64 // tentative distance, sum of arc weights
	nx   int     // slice index, "node id"
	fx   int     // heap.Fix index
	done bool
}

type tent []*tentResult

// DijkstraPath finds a single shortest path.
//
// Returned is the path and distance as returned by WeightedFromTree.PathTo.
func DijkstraPath(g WeightedAdjacencyList, start, end int) ([]Half, float64) {
	d := NewDijkstra(g)
	d.Path(start, end)
	return d.Result.PathTo(end)
}

// Path finds a single shortest path.
//
// Path returns true if a path exists, false if not. The path can be recovered
// from b.Result.  Path returns as soon as the shortest path to end is found;
// it does not explore the remainder of the graph.
func (d *Dijkstra) Path(start, end int) bool {
	d.search(start, end)
	return d.Result.Paths[end].Len > 0
}

// AllPaths finds shortest paths from start to all nodes reachable
// from start.
//
// AllPaths returns number of paths found, equivalent to the number of nodes
// reached, including the path ending at start.  Path results are left in
// d.Result.
func (d *Dijkstra) AllPaths(start int) (nFound int) {
	return d.search(start, -1)
}

// returns number of nodes reached (= number of shortest paths found)
func (d *Dijkstra) search(start, end int) (reached int) {
	// reset from any previous run
	d.ndVis = 0
	d.arcVis = 0
	d.Result.reset()
	for i := range d.r {
		d.r[i].done = false
	}

	current := start
	rp := d.Result.Paths
	rp[current].Len = 1 // path length at start is 1 node
	cr := &d.r[current]
	cr.dist = 0 // distance at start is 0
	rp[current].Dist = 0
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
			dist := cr.dist + nb.ArcWeight
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
			rp[nb.To].From = FromHalf{current, nb.ArcWeight}
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
		rp[current].Dist = cr.dist // store final distance
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
func AStarAPath(g WeightedAdjacencyList, start, end int, h Heuristic) ([]Half, float64) {
	a := NewAStar(g)
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
func AStarMPath(g WeightedAdjacencyList, start, end int, h Heuristic) ([]Half, float64) {
	a := NewAStar(g)
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

// A BellmanFord object allows shortest path searches using the
// Bellman-Ford-Moore algorithm.
type BellmanFord struct {
	Graph  WeightedAdjacencyList
	Result *WeightedFromTree
}

// NewBellmanFord creates a BellmanFord object that allows shortest path
// searches using the Bellman-Ford-Moore algorithm.
//
// Argument g is the graph to be searched, as a weighted adjacency list.
// Negative arc weights are allowed as long as there are no negative cycles.
// Graphs may be directed or undirected.  Loops and parallel arcs are
// allowed.
//
// The graph g will not be modified by any BellmanFord methods.  NewBellmanFord
// initializes the BellmanFord object for the length (number of nodes) of g.
// If you add nodes to your graph, abandon any previously created BellmanFord
// object and call NewBellmanFord again.
//
// NewBellmanFord calls NewWeightedFromTree.  See documentation there in
// particular for the option to change NoPath before running a search.
//
// Searches on a single BellmanFord object can be run consecutively but not
// concurrently.  Searches can be run concurrently however, on BellmanFord
// objects obtained with separate calls to NewBellmanFord, even with the same
// graph argument to NewBellmanFord.
func NewBellmanFord(g WeightedAdjacencyList) *BellmanFord {
	return &BellmanFord{g, NewWeightedFromTree(len(g))}
}

// Run runs the BellmanFord algorithm, which finds shortest paths from start
// to all nodes reachable from start.
//
// The algorithm allows negative edge weights but not negative cycles.
// Run returns true if the algorithm completes successfully.  In this case
// b.Result will be populated with a WeightedFromTree encoding shortest paths
// found.
//
// Run returns false in the case that it encounters a negative cycle.
// In this case values in b.Result are meaningless.
func (b *BellmanFord) Run(start int) (ok bool) {
	b.Result.reset()
	rp := b.Result.Paths
	rp[start].Dist = 0
	rp[start].Len = 1
	for _ = range b.Graph[1:] {
		imp := false
		for from, nbs := range b.Graph {
			fp := &rp[from]
			d1 := fp.Dist
			for _, nb := range nbs {
				d2 := d1 + nb.ArcWeight
				to := &rp[nb.To]
				// TODO improve to break ties
				if fp.Len > 0 && d2 < to.Dist {
					*to = WeightedPathEnd{
						Dist: d2,
						From: FromHalf{from, nb.ArcWeight},
						Len:  fp.Len + 1,
					}
					imp = true
				}
			}
		}
		if !imp {
			break
		}
	}
	for from, nbs := range b.Graph {
		d1 := rp[from].Dist
		for _, nb := range nbs {
			if d1+nb.ArcWeight < rp[nb.To].Dist {
				return false // negative cycle
			}
		}
	}
	return true
}
