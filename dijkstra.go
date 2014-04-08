// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"container/heap"
	"math"
)

// A Dijkstra object allows shortest path searches using Dijkstra's algorithm.
//
// Construct with NewDijkstra.  NoPath is used as a distance result when no
// path or no arc exists.  It is initialized to +Inf by NewDijkstra but you
// can assign a different valule to this field if you like.
type Dijkstra struct {
	Graph  WeightedAdjacencyList
	NoPath float64 // initialized to +Inf by NewDijkstra.
	Start  int
	Result []DijkstraResult
	MaxLen int
	// test instrumentation
	ndVis, arcVis int
}

// NewDijkstra creates a Dijkstra struct that allows shortest path searches.
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
// Searches on a single Dijkstra object can be run consecutively but not
// concurrently.  Searches can be run concurrently however, on Dijkstra
// objects obtained with separate calls to NewDijkstra, even with the same
// graph argument to NewDijkstra.
func NewDijkstra(g WeightedAdjacencyList) *Dijkstra {
	r := make([]DijkstraResult, len(g))
	for i := range r {
		r[i].nx = i
	}
	return &Dijkstra{
		Graph:  g,
		Result: r,
		NoPath: math.Inf(1),
	}
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
type DijkstraResult struct {
	PathDist float64  // sum of arc weights
	FromTree HalfFrom // encodes shortest path
	PathLen  int      // number of nodes

	nx   int // slice index, "node id"
	fx   int // heap.Fix index
	done bool
}

type tent []*DijkstraResult

// Path finds a single shortest path from start to end.
//
// Returned is the path and distance as returned by Dijkstra.PathTo.
// The function returns as soon as the shortest path to end is found.
// It does not explore the remainder of the graph.
func (d *Dijkstra) Path(start, end int) ([]Half, float64) {
	d.search(start, end)
	return d.PathTo(end)
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
func (d *Dijkstra) PathTo(end int) ([]Half, float64) {
	n := d.Result[end].PathLen
	if n == 0 {
		return nil, d.NoPath
	}
	p := make([]Half, n)
	dist := 0.
	for {
		n--
		f := &d.Result[end].FromTree
		p[n] = Half{end, f.ArcWeight}
		if n == 0 {
			return p, dist
		}
		dist += f.ArcWeight
		end = f.From
	}
}

// AllPaths finds shortest paths from start to all nodes reachable
// from start.  Results are in d.Result; individual paths can be decoded
// with DijkstraResult.PathTo.
//
// Returned is the number of nodes reached, or the number of paths found,
// including the path ending at start.
func (d *Dijkstra) AllPaths(start int) int {
	return d.search(start, -1)
}

// returns number of nodes reached (= number of shortest paths found)
func (d *Dijkstra) search(start, end int) int {
	// reset from any previous run
	d.ndVis = 0
	d.arcVis = 0
	for i := range d.Result {
		r := &d.Result[i]
		r.done = false
		r.PathLen = 0
		r.PathDist = d.NoPath
		r.FromTree = HalfFrom{-1, d.NoPath}
	}

	current := start
	cr := &d.Result[current]
	cr.PathLen = 1  // path length 1 for start node
	cr.PathDist = 0 // distance at start is 0
	cr.done = true  // mark start done.  it skips the heap.
	nDone := 1      // accumulated for a return value
	var t tent
	for current != end {
		for _, nb := range d.Graph[current] {
			d.arcVis++
			hr := &d.Result[nb.To]
			if hr.done {
				continue // skip nodes already done
			}
			dist := cr.PathDist + nb.ArcWeight
			visited := hr.PathLen > 0
			if visited && dist >= hr.PathDist {
				continue // it's no better
			}
			// the path through current to this node is shortest so far.
			// record new path data for this node and update tentative set.
			hr.PathDist = dist
			hr.PathLen = cr.PathLen + 1
			hr.FromTree = HalfFrom{current, nb.ArcWeight}
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
		cr = heap.Pop(&t).(*DijkstraResult)
		cr.done = true
		nDone++
		current = cr.nx
	}
	return -1 // normal return for single shortest path search
}

// tent implements container/heap
func (t tent) Len() int           { return len(t) }
func (t tent) Less(i, j int) bool { return t[i].PathDist < t[j].PathDist }
func (t tent) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].fx = i
	t[j].fx = j
}
func (s *tent) Push(x interface{}) {
	nd := x.(*DijkstraResult)
	nd.fx = len(*s)
	*s = append(*s, nd)
}
func (s *tent) Pop() interface{} {
	t := *s
	last := len(t) - 1
	*s = t[:last]
	return t[last]
}
