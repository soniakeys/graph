// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Ed implements Dijkstra's shortest path algorithm.
//
// This is a stripped down version of github.com/soniakeys/graph.  It serves
// a couple of purposes.  One purpose is understanding the algorithm.
// This package shows it in a slightly simpler form that might be a little
// easier to figure out.  Another purpose is exploring the overhead of making
// the package reusable, memory efficient, and goroutine safe.
//
// Graph.DijkstraShortest path is designed to reusable by operating through
// interfaces rather than directly on structs.  A user of the the package has
// to implement interfaces of the package but does not have to adapt the code
// of the package.  This is best for usability but has some small runtime cost.
//
// A greater cost in graph comes from its minimal use of memory.  The algorithm
// requires some bookkeeping information per node, but not all data is needed
// for all nodes at all times.  Graph attempts to allocate data only as
// needed, but there is overhead involved in maintaining and traversing
// these extra pointers.
//
// Finally, graph is goroutine safe.  Multiple searches can be run
// concurrently on a single graph.  For this, the algorithm must maintain
// the bookkeeping data per search.  To do this without bothering the user
// of the package requires associating the user's nodes with the bookkeeping
// data.  This is done with a Go map.
//
// This package explores an alternative, a simple network of linked structs,
// no interfaces, no maps.  It turns out to run several times faster.
// That sounds significant but surely for many applications it wouldn't matter.
// Much more significant would seem to be the difference in coding effort.
// To use this package you can't simply read the docs and code to an API,
// you must edit the code and adapt it to your application.
// If your application needs concurrent searches on the same graph, you
// are on your own to solve the problem.
package ed

import "container/heap"

// Node implements a node or vertex in the graph.  You must use this
// defintion as a template and adapt the application specific data and
// perhaps the graph representation to your application.
type Node struct {
	Name string // application specific

	Neighbors []Neighbor // graph representation

	// fields used for nodes visited in shortest path computation
	done     bool
	prevNode *Node   // path back to start
	prevEdge float64 // distance from prev node to the node of this struct
	// fields used for nodes in the tentative set
	dist float64 // tentative path distance
	n    int     // number of nodes in path
	rx   int     // heap.Remove index
}

// String satisfies fmt.Stringer
func (n *Node) String() string {
	return n.Name
}

// Reset prepares a node to be searched again.  The graph can be searched
// once without preparation, but before another search you must traverse
// the graph and call this function on each node.  Note there is no predefined
// mechanism for this.  You must implement some collection that holds all
// nodes in the graph and allows iteration over them.  (A slice or map
// should do.)
func (n *Node) Reset() {
	n.n = 0
	n.done = false
}

// Neighbor represents a directed edge.  It can be adapted as needed for
// your application.
type Neighbor struct {
	Edge float64
	Node *Node
}

type tent []*Node

// ShortestPath runs Dijkstra's shortest path algorithm, returning the
// shortest path from start to end.  Bookkeeping data is stored in the
// nodes and so only a singe search can be run at a time.  Before a
// subsequent search, call Reset as described in the documentation for
// that function.
func ShortestPath(start, end *Node) ([]Neighbor, float64) {
	var t tent // tentative set, heap
	current := start
	current.n = 1       // path length 1 for start node
	current.done = true // mark start done.  it skips the heap.
	for {
		for _, nb := range current.Neighbors {
			nd := nb.Node
			if nd.done {
				continue // skip nodes already done
			}
			dist := current.dist + nb.Edge
			visited := nd.n > 0
			if visited && dist >= nd.dist {
				continue // it's no better
			}
			// the path through current to this node is shortest so far.
			// record new path data for this node and update tentative set.
			nd.dist = dist
			nd.n = current.n + 1
			nd.prevNode = current
			nd.prevEdge = nb.Edge
			if visited {
				heap.Fix(&t, nd.rx)
			} else {
				heap.Push(&t, nd)
			}
		}
		if current == end { // search complete
			// recover path by tracing prev links
			i := current.n
			path := make([]Neighbor, i)
			for n := current; n != nil; {
				i--
				path[i] = Neighbor{n.prevEdge, n}
				n = n.prevNode
			}
			return path, current.dist // success
		}
		if len(t) == 0 {
			return nil, 0 // failure. no more reachable nodes
		}
		// new current is node with smallest tentative distance
		current = heap.Pop(&t).(*Node)
		current.done = true
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
	r := t[last]
	*s = t[:last]
	return r
}
