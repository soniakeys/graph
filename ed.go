// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Ed implements Dijkstra's single shortest path algorithm.
//
// The data representation is fixed and holds no more information than what
// is needed to run the search algorithm.  To search your own data you must
// create a graph using the types of this package, run the search, and use
// the result to navigate your data.  The types here use integer node IDs.
// If your data uses something else you must devise a translation to integer
// node IDs.
package ed

import (
	"container/heap"
)

type Dijkstra struct {
	g   [][]Half
	dat []ndDat
	// instrumentation
	ndVis, arcVis int
}

// New creates a Dijkstra struct that allows shortest path searches.
//
// Argument g is the graph to be searched, as an adjacency list where node
// IDs correspond to the slice indexes of g.  Each []Half element of g
// represents the neighbors of a node.  To represent an undirected graph,
// create reciprocal Halfs with identical ArcWeights.
// Nodes cannot be added or removed to a Dijkstra object.
func New(g [][]Half) *Dijkstra {
	dat := make([]ndDat, len(g))
	for i := range dat {
		dat[i].nx = i
	}
	return &Dijkstra{g: g, dat: dat}
}

// Half is a half arc, representing a "neighbor" of a node.
//
// Halfs can be composed to form an adjacency list.
type Half struct {
	To        int // index in graph slice
	ArcWeight float64
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
// the single shortest path from start to end.  Searches can be run
// consecutively but not concurrently.
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
	if len(t) == 0 {
		return nil
	}
	last := len(t) - 1
	*s = t[:last]
	return t[last]
}
