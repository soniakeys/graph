// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"container/heap"
)

// A Dijkstra object allows shortest path searches using Dijkstra's algorithm.
type Dijkstra struct {
	g   [][]Half // graph supplied by user.  this is not modified.
	dat []ndDat  // working data for the algorithm
	// instrumentation
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
func New(g [][]Half) *Dijkstra {
	dat := make([]ndDat, len(g))
	for i := range dat {
		dat[i].nx = i
	}
	return &Dijkstra{g: g, dat: dat}
}

// ndDat. per node bookeeping data needed for Dijktra's algorithm.
type ndDat struct {
	nx int // index in graph slice, "node id"
	// fields used for nodes visited in shortest path computation
	done       bool
	prevFrom   int     // path back to start
	prevWeight float64 // weight of arc from prev node
	// fields used for nodes in the tentative set
	dist float64 // tentative path distance
	n    int     // number of nodes in path
	rx   int     // heap.Remove index
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
	// reset from any previous run
	d.ndVis = 0
	d.arcVis = 0
	for i := range d.dat {
		d.dat[i].n = 0
		d.dat[i].done = false
	}

	current := start
	cn := &d.dat[current]
	cn.n = 1       // path length 1 for start node
	cn.done = true // mark start done.  it skips the heap.
	var t tent
	for {
		for _, nb := range d.g[current] {
			d.arcVis++
			if nb.To == end {
				// search complete
				// recover path by tracing prev links
				i := cn.n
				dist := cn.dist + nb.ArcWeight
				path := make([]Half, i+1)
				path[i] = nb
				for n := current; i > 0; n = cn.prevFrom {
					cn = &d.dat[n]
					i--
					path[i] = Half{n, cn.prevWeight}
				}
				return path, dist // success
			}
			hn := &d.dat[nb.To]
			if hn.done {
				continue // skip nodes already done
			}
			dist := cn.dist + nb.ArcWeight
			visited := hn.n > 0
			if visited && dist >= hn.dist {
				continue // it's no better
			}
			// the path through current to this node is shortest so far.
			// record new path data for this node and update tentative set.
			hn.dist = dist
			hn.n = cn.n + 1
			hn.prevFrom = current
			hn.prevWeight = nb.ArcWeight
			if visited {
				heap.Fix(&t, hn.rx)
			} else {
				heap.Push(&t, hn)
			}
		}
		d.ndVis++
		if len(t) == 0 {
			return nil, 0 // failure. no more reachable nodes
		}
		// new current is node with smallest tentative distance
		cn = heap.Pop(&t).(*ndDat)
		cn.done = true
		current = cn.nx
	}
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
