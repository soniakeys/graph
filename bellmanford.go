// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

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
