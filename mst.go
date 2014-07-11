// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"container/heap"
	"log"
)

// Prim implements the Jarník-Prim-Dijkstra algorithm for constructing
// a minimum spanning tree on an undirected graph.
//
// Construct with NewPrim.
type Prim struct {
	Graph  WeightedAdjacencyList
	Result *WeightedFromTree

	best []prNode // slice backs heap
}

// NewPrim constructs a new Prim object.  Argument g should represent an
// undirected graph.
func NewPrim(g WeightedAdjacencyList) *Prim {
	b := make([]prNode, len(g))
	for n := range b {
		b[n].nx = n
		b[n].fx = -1
	}
	return &Prim{
		Graph:  g,
		Result: NewWeightedFromTree(len(g)),
		best:   b,
	}
}

type prNode struct {
	nx   int
	from FromHalf
	fx   int
}

type prHeap []*prNode

// Reset clears results of Span, allowing results to be recomputed.
//
// Reset is not meaningful following a change to the number of nodes in the
// graph.  To recompute following addition or deletion of nodes, simply
// abandon the Prim object and create a new one.
func (p *Prim) Reset() {
	p.Result.reset()
	b := p.best
	for n := range b {
		b[n].fx = -1
	}
}

// Span computes a minimal spanning tree on the connected component containing
// the given start node.
//
// If a graph has multiple connected components, a spanning forest can be
// accumulated by calling Span successively on representative nodes of the
// components.  Note that Result.Start can contain only the most recent start
// node.  To complete the forest representation you must retain a separate
// list of start nodes.
func (p *Prim) Span(start int) int {
	rp := p.Result.Paths
	var frontier prHeap
	rp[start].Dist = 0
	rp[start].Len = 1
	b := p.best
	nDone := 1
	for a := start; ; {
		log.Println("node", a)
		for _, nb := range p.Graph[a] {
			log.Println("  nb", nb)
			if rp[nb.To].Len > 0 {
				continue // already in MST, no action
			}
			switch bp := &b[nb.To]; {
			case bp.fx == -1: // new node for frontier
				bp.from = FromHalf{a, nb.ArcWeight}
				heap.Push(&frontier, bp)
			case nb.ArcWeight < bp.from.ArcWeight: // better arc
				bp.from = FromHalf{a, nb.ArcWeight}
				log.Println("heap")
				for _, pn := range frontier {
					log.Printf("  %#v\n", pn)
				}
				heap.Fix(&frontier, bp.fx)
			}
		}
		if len(frontier) == 0 {
			break // done
		}
		bp := heap.Pop(&frontier).(*prNode)
		a = bp.nx
		rp[a].Len = rp[bp.from.From].Len + 1
		rp[a].From = bp.from
		nDone++
	}
	return nDone
}

func (h prHeap) Len() int { return len(h) }
func (h prHeap) Less(i, j int) bool {
	return h[i].from.ArcWeight < h[j].from.ArcWeight
}
func (h prHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].fx = i
	h[j].fx = j
}
func (p *prHeap) Push(x interface{}) {
	nd := x.(*prNode)
	nd.fx = len(*p)
	*p = append(*p, nd)
}
func (p *prHeap) Pop() interface{} {
	r := *p
	last := len(r) - 1
	*p = r[:last]
	return r[last]
}
