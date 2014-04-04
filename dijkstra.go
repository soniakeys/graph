// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"container/heap"
	"math"
)

// A Dijkstra object allows shortest path searches using Dijkstra's algorithm.
type Dijkstra struct {
	Graph  [][]Half
	Result []DijkstraResult
	// test instrumentation
	ndVis, arcVis int
}

// NewDijkstra creates a Dijkstra struct that allows shortest path searches.
//
// Argument g is the graph to be searched, as an adjacency list where node
// IDs correspond to the slice indexes of g.  Each []Half element of g
// represents the neighbors of a node.  All Half.To fields must contain
// a node ID greater than or equal to zero and strictly less than len(g).
// As usual for Dijkstra's algorithm, arc weights must be non-negative.
//
// The algorithm conceptually works for directed and undirected graphs but
// the adjacency list structure is inherently directed.  To represent an
// undirected graph, create reciprocal Halfs with identical arc weights.
//
// Loops (arcs from a node to itself) are allowed but will not affect the
// result.  Parallel arcs (multiple arcs between the same two nodes) are
// also allowed; the algorithm will find the shortest one.  Generally though,
// you should create simple graphs (graphs with no loops or parallel arcs)
// where convenient, for best performance.
//
// The graph g will not be modified by any Dijkstra methods.  New initializes
// the Dijkstra object for the length (number of nodes) of g.  If you add
// nodes to your graph, abandon any previously created Dijkstra object and
// call New again.
func NewDijkstra(g [][]Half) *Dijkstra {
	r := make([]DijkstraResult, len(g))
	for i := range dat {
		r[i].nx = i
	}
	return &Dijkstra{
		Graph:  g,
		Result: r,
	}
}

type DijkstraResult struct {
	PathLen  int
	PathDist float64
	FromTree FromHalf
	nx       int // index in graph slice, "node id"
	fx       int // heap.Fix index
	done     bool
}

type tent []*ndDat

// instrumentation
func (d *Dijkstra) na() (int, int) {
	return d.ndVis, d.arcVis
}

// SingleShortestPath runs Dijkstra's shortest path algorithm, returning
// the single shortest path from start to end.
//
// Searches on a single Dijkstra object can be run consecutively but not
// concurrently.  Searches can be run concurrently however, on Dijkstra
// objects obtained with separate calls to New, even with the same graph
// argument to New.
//
// The slice result represents the found path with a sequence of half arcs.
// For the first element, representing the start node, the arc length is
// meaningless and will always be 0.  If no path exists from start to end
// the slice result will be nil.  The total path length, as the sum of the
// arc lengths, is also returned.
func (d *Dijkstra) SingleShortestPath(start, end int) ([]Half, float64) {
	if start == end {
		return []Half{{end, 0}}, 0
	}
	d.search(start, end)
	return d.PathTo(end)
}

func (d *Dijkstra) PathTo(end int) ([]Half, float64) {
	n := d.PathLen[end]
	if n == 0 {
		return nil, math.Inf(1)
	}
	p := make([]Half, n)
	dist := 0.
	for n--; n >= 0; n-- {
		f := d.FromTree[end]
		p[n] = Half{end, f.ArcWeight}
		dist += f.ArcWeight
		end = f.From
	}
	return p, dist
}

func (d *Dijkstra) AllShortestPaths(start int) int {
	return d.search(start, -1)
}

// returns number of nodes reached (= number of shortest paths found)
func (d *Dijkstra) search(start, end int) int {
	// reset from any previous run
	d.ndVis = 0
	d.arcVis = 0
	for i := range d.dat {
		d.PathLen[i] = 0
		d.PathDist[i] = 0
		d.dat[i].done = false
		d.FromTree[i] = FromHalf{-1, 0}
	}

	current := start
	d.PathLen[current] = 1 // path length 1 for start node
	cn := &d.dat[current]
	cn.done = true // mark start done.  it skips the heap.
	nDone := 1
	var t tent
	for current != end {
		for _, nb := range d.Graph[current] {
			d.arcVis++
			hn := &d.dat[nb.To]
			if hn.done {
				continue // skip nodes already done
			}
			dist := d.PathDist[current] + nb.ArcWeight
			visited := d.PathLen[nb.To] > 0
			if visited && dist >= d.PathDist[nb.To] {
				continue // it's no better
			}
			// the path through current to this node is shortest so far.
			// record new path data for this node and update tentative set.
			d.PathDist[nb.To] = dist
			d.PathLen[nb.To] = d.PathLen[current] + 1
			d.FromTree[nb.To] = FromHalf{current, nb.ArcWeight}
			if visited {
				heap.Fix(&t, hn.rx)
			} else {
				heap.Push(&t, hn)
			}
		}
		d.ndVis++
		if len(t) == 0 {
			return nDone // no more reachable nodes. AllPaths normal return
		}
		// new current is node with smallest tentative distance
		cn = heap.Pop(&t).(*ndDat)
		cn.done = true
		nDone++
		current = cn.nx
	}
	return -1 // normal return for single shortest path search
}

// tent implements container/heap
func (t tent) Len() int           { return len(t) }
func (t tent) Less(i, j int) bool { return t[i].dist < t[j].dist }
func (t tent) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].rx = i
	t[j].rx = j
}
func (s *tent) Push(x interface{}) {
	nd := x.(*ndDat)
	nd.rx = len(*s)
	*s = append(*s, nd)
}
func (s *tent) Pop() interface{} {
	t := *s
	last := len(t) - 1
	*s = t[:last]
	return t[last]
}
