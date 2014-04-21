// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"container/heap"
)

type Prim struct {
	Graph  WeightedAdjacencyList
	Result *WeightedFromTree

	best []prNode // slice backs heap
}

func NewPrim(g WeightedAdjacencyList) *Prim {
	b := make([]prNode, len(g))
	for n := range b {
		b[n].nx = n
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

func (p *Prim) Run(start int) int {
	p.Result.reset()
	b := p.best
	for n := range b {
		b[n].fx = -1
	}
	rp := p.Result.Paths
	var frontier prHeap
	rp[start].Dist = 0
	rp[start].Len = 1
	nDone := 1
	for a := start; ; {
		for _, nb := range p.Graph[a] {
			if rp[nb.To].Len > 0 {
				continue // already in MST, no action
			}
			switch bp := &b[nb.To]; {
			case bp.fx == -1: // new node for frontier
				bp.from = FromHalf{a, nb.ArcWeight}
				heap.Push(&frontier, bp)
			case nb.ArcWeight < bp.from.ArcWeight: // better arc
				bp.from = FromHalf{a, nb.ArcWeight}
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
func (h prHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
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
