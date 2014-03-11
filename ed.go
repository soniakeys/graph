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

// Graph reprsentation is a slice of nodes, linked to each other by adjacency
// lists.
type Graph []Node

// Half is a half arc, representing a "neighbor" of a node.
type Half struct {
	To        int // index in graph slice
	ArcWeight float64
}

// Each Nodes has an adjacency list and bookeeping data needed for Dijktra's
// algorithm.
type Node struct {
	Nbs []Half // adjacency list, "neighbors"
	nx  int    // index in graph slice, "node id"
	// fields used for nodes visited in shortest path computation
	done       bool
	prevFrom   int     // path back to start
	prevWeight float64 // weight of arc from prev node
	// fields used for nodes in the tentative set
	dist float64 // tentative path distance
	n    int     // number of nodes in path
	rx   int     // heap.Remove index
}

var ndVis, arcVis int // instrumentation

type tent []*Node

// New constructs and initializes a graph with n nodes.
func New(n int) Graph {
	g := make(Graph, n)
	for i := range g {
		g[i].nx = i
	}
	return g
}

// SetArcs sets the adjacency list for a node, replacing the existing list.
// SetArcs panics if argument from is not a valid index for the nodes of the
// graph.  No validation is performed on argument to.
func (g Graph) SetArcs(from int, to []Half) {
	g[from].Nbs = to
}

// AddArcSimple adds an arc to a graph safely, validating the arguments and
// maintaining the the graph as a "simple" graph, that is, with no self-loops
// and with no multiple arcs from one node to another.
//
// AddArcSimple returns true if 1) node indexes from and to.To are valid indexes
// for the graph and not equal to each other, 2) if to.ArcWeight is non-negative
// (and non-NaN), and 3) if the the arc is not already in the graph.  Otherwise
// the graph is not modified and the function returns false.
func (g Graph) AddArcSimple(from int, to Half) bool {
	if from < 0 || from >= len(g) || to.To < 0 || to.To >= len(g) {
		return false
	}
	if !(to.ArcWeight >= 0) { // inverse test to catch NaNs
		return false
	}
	nbs := &g[from].Nbs
	if from == to.To {
		return false // disallow loops
	}
	for _, nb := range *nbs {
		if nb.To == to.To {
			return false // disallow parallel arcs
		}
	}
	*nbs = append(*nbs, to)
	return true
}

// ShortestPath runs Dijkstra's shortest path algorithm, returning the single
// shortest path from start to end.  Searches can be run consecutively but not
// concurrently.
func (g Graph) ShortestPath(start, end int) ([]Half, float64) {
	//	fmt.Println("start", start, "end", end)
	if start == end {
		return []Half{{end, 0}}, 0
	}
	// reset g from any previous run
	ndVis = 0
	arcVis = 0
	for i := range g {
		g[i].n = 0
		g[i].done = false
	}

	current := start
	cn := &g[current]
	cn.n = 1       // path length 1 for start node
	cn.done = true // mark start done.  it skips the heap.
	var t tent
	for {
		//		fmt.Println("current", current)
		for _, nb := range cn.Nbs {
			//			fmt.Printf("  nb %#v\n", nb)
			arcVis++
			if nb.To == end {
				// search complete
				// recover path by tracing prev links
				i := cn.n
				dist := cn.dist + nb.ArcWeight
				path := make([]Half, i+1)
				path[i] = nb
				for n := current; i > 0; n = cn.prevFrom {
					cn = &g[n]
					i--
					path[i] = Half{n, cn.prevWeight}
				}
				return path, dist // success
			}
			hn := &g[nb.To]
			if hn.done {
				//				fmt.Println("    done")
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
		ndVis++
		if len(t) == 0 {
			return nil, 0 // failure. no more reachable nodes
		}
		// new current is node with smallest tentative distance
		cn = heap.Pop(&t).(*Node)
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
	nd := x.(*Node)
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
