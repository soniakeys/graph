// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"container/heap"
)

// A Dijkstra object allows shortest path searches using Dijkstra's algorithm.
//
// Construct with NewDijkstra.  NoPath is used as a distance result when no
// path or no arc exists.  It is initialized to +Inf by NewDijkstra but you
// can assign a different valule to this field if you like.
type Dijkstra struct {
	Graph  WeightedAdjacencyList
	Result *WeightedFromTree
	r      []tentResult
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

// Path finds a single shortest path from start to end.
//
// Returned is the path and distance as returned by Dijkstra.PathTo.
// The function returns as soon as the shortest path to end is found.
// It does not explore the remainder of the graph.
func (d *Dijkstra) Path(start, end int) ([]Half, float64) {
	d.search(start, end)
	return d.Result.PathTo(end)
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
func (d *Dijkstra) search(start, end int) (reached int) {
	// reset from any previous run
	d.ndVis = 0
	d.arcVis = 0
	d.Result.Reset()
	for i := range d.r {
		r := &d.r[i]
		r.done = false
		r.dist = d.Result.NoPath
	}

	current := start
	rp := d.Result.Paths
	rp[current].Len = 1 // path length 1 for start node
	cr := &d.r[current]
	cr.dist = 0    // distance at start is 0
	rp[current].Dist = 0
	cr.done = true // mark start done.  it skips the heap.
	nDone := 1     // accumulated for a return value
	var t tent
	for current != end {
		for _, nb := range d.Graph[current] {
			d.arcVis++
			hr := &d.r[nb.To]
			if hr.done {
				continue // skip nodes already done
			}
			dist := cr.dist + nb.ArcWeight
			visited := rp[nb.To].Len > 0
			if visited && dist >= hr.dist {
				continue // it's no better
			}
			// the path through current to this node is shortest so far.
			// record new path data for this node and update tentative set.
			hr.dist = dist
			rp[nb.To].Len = rp[current].Len + 1
			rp[nb.To].From = HalfFrom{current, nb.ArcWeight}
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
