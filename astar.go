// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"container/heap"
	"math"
)

type AStar struct {
	g [][]Half // input graph
	// r is a list of all nodes reached so far.
	// the chain of nodes following the prev member represents the
	// best path found so far from the start to this node.
	r []rNode
}

func NewAStar(g [][]Half) *AStar {
	r := make([]rNode, len(g))
	for i := range r {
		r[i].nd = i
	}
	return &AStar{g: g, r: r}
}

// rNode holds data for a "reached" node
type rNode struct {
	nd       int
	state    int8    // state constants defined below
	prevNode *rNode  // chain encodes path back to start
	prevEdge float64 // edge from prevNode to the node of this struct
	g        float64 // "g" best known path distance from start node
	f        float64 // "g+h", path dist + heuristic estimate
	n        int     // number of nodes in path
	rx       int     // heap.Remove index
}

// for rNode.state
const (
	unreached = 0
	reached   = 1
	open      = 1
	closed    = 2
)

type openHeap []*rNode

func (a *AStar) AStarA(start, end int, h func(int) float64) ([]Half, float64) {
	// start node is reached initially
	p := &a.r[start]
	p.state = reached
	p.f = h(start) // total path estimate is estimate from start
	p.n = 1        // path length is 1 node
	// oh is a heap of nodes "open" for exploration.  nodes go on the heap
	// when they get an initial or new "g" path distance, and therefore a
	// new "f" which serves as priority for exploration.
	oh := openHeap{p}
	for len(oh) > 0 {
		bestPath := heap.Pop(&oh).(*rNode)
		bestNode := bestPath.nd
		if bestNode == end {
			// done
			dist := bestPath.g
			i := bestPath.n
			path := make([]Half, i)
			for i > 0 {
				i--
				path[i] = Half{To: bestPath.nd, ArcWeight: bestPath.prevEdge}
				bestPath = bestPath.prevNode
			}
			return path, dist
		}
		for _, nb := range a.g[bestNode] {
			ed := nb.ArcWeight
			nd := nb.To
			g := bestPath.g + ed
			if alt := &a.r[nd]; alt.state == reached {
				if g > alt.g {
					// new path to nd is longer than some alternate path
					continue
				}
				if g == alt.g && bestPath.n+1 >= alt.n {
					// new path has identical length of some alternate path
					// but it takes more hops.  stick with fewer nodes in path.
					continue
				}
				// cool, we found a better way to get to this node.
				// update alt with new data and make sure it's on the heap.
				alt.prevNode = bestPath
				alt.prevEdge = ed
				alt.g = g
				alt.f = g + h(nd)
				alt.n = bestPath.n + 1
				if alt.rx < 0 {
					heap.Push(&oh, alt)
				} else {
					heap.Fix(&oh, alt.rx)
				}
			} else {
				// bestNode being reached for the first time.
				alt.state = reached
				alt.prevNode = bestPath
				alt.prevEdge = ed
				alt.g = g
				alt.f = g + h(nd)
				alt.n = bestPath.n + 1
				heap.Push(&oh, alt) // and it's now open for exploration
			}
		}
	}
	return nil, math.Inf(1) // no path
}

func (a *AStar) AStarM(start, end int, h func(int) float64) ([]Half, float64) {
	p := &a.r[start]
	p.f = h(start) // total path estimate is estimate from start
	p.n = 1        // path length is 1 node

    // difference from AStarA:
    // instead of a bit to mark a reached node, there are two states,
    // open and closed. open marks nodes "open" for exploration.
    // nodes are marked open as they are reached, then marked
    // closed as they are found to be on the best path.
	p.state = open

	oh := openHeap{p}
	for len(oh) > 0 {
		bestPath := heap.Pop(&oh).(*rNode)
		bestNode := bestPath.nd
		if bestNode == end {
			// done
			dist := bestPath.g
			i := bestPath.n
			path := make([]Half, i)
			for i > 0 {
				i--
				path[i] = Half{To: bestPath.nd, ArcWeight: bestPath.prevEdge}
				bestPath = bestPath.prevNode
			}
			return path, dist
		}

        // difference from AStarA:
        // move nodes to closed list as they are found to be best so far.
		bestPath.state = closed

		for _, nb := range a.g[bestNode] {
			ed := nb.ArcWeight
			nd := nb.To

            // difference from AStarA:
            // Monotonicity means that f cannot be improved.
            if a.r[nd].state == closed {
                continue
            }

			g := bestPath.g + ed
			if alt := &a.r[nd]; alt.state == open {
				if g > alt.g {
					// new path to nd is longer than some alternate path
					continue
				}
				if g == alt.g && bestPath.n+1 >= alt.n {
					// new path has identical length of some alternate path
					// but it takes more hops.  stick with fewer nodes in path.
					continue
				}
				// cool, we found a better way to get to this node.
				// update alt with new data and make sure it's on the heap.
				alt.prevNode = bestPath
				alt.prevEdge = ed
				alt.g = g
				alt.f = g + h(nd)
				alt.n = bestPath.n + 1

                // difference from AStarA:
                // we know alt was on the heap because we found it marked open
				heap.Fix(&oh, alt.rx)
			} else {
				// bestNode being reached for the first time.
				alt.state = open
				alt.prevNode = bestPath
				alt.prevEdge = ed
				alt.g = g
				alt.f = g + h(nd)
				alt.n = bestPath.n + 1
				heap.Push(&oh, alt) // and it's now open for exploration
			}
		}
	}
	return nil, math.Inf(1) // no path
}

// implement container/heap
func (h openHeap) Len() int           { return len(h) }
func (h openHeap) Less(i, j int) bool { return h[i].f < h[j].f }
func (h openHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].rx = i
	h[j].rx = j
}
func (p *openHeap) Push(x interface{}) {
	h := *p
	rx := len(h)
	h = append(h, x.(*rNode))
	h[rx].rx = rx
	*p = h
}

func (p *openHeap) Pop() interface{} {
	h := *p
	last := len(h) - 1
	*p = h[:last]
	h[last].rx = -1
	return h[last]
}
