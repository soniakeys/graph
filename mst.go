// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"container/heap"
	"sort"
)

type dsElement struct {
	from int
	rank int
}

type disjointSet struct {
	set []dsElement
}

func newDisjointSet(n int) disjointSet {
	set := make([]dsElement, n)
	for i := range set {
		set[i].from = -1
	}
	return disjointSet{set}
}

func (ds disjointSet) find(n int) int {
	// fast paths for n == root or from root.
	// no updates need in these cases.
	s := ds.set
	fr := s[n].from
	if fr < 0 { // n is root
		return n
	}
	n, fr = fr, s[fr].from
	if fr < 0 { // n is from root
		return n
	}
	// otherwise updates needed.
	// two iterative passes (rather than recursion or stack)
	// pass 1: find root
	r := fr
	for {
		f := s[r].from
		if f < 0 {
			break
		}
		r = f
	}
	// pass 2: update froms
	for {
		s[n].from = r
		if fr == r {
			return r
		}
		n = fr
		fr = s[n].from
	}
}

// Kruskal implements Kruskal's algorithm for constructing a minimum spanning
// tree on an undirected graph.
//
// The tree is returned as a FromList, with leaf nodes indicated, but without
// path lengths or MaxPathLen set.  Also returned is a parallel list of labels
// and a total distance for the returned tree.
//
// The edge list of the receiver is sorted as a side effect of this method.
// See KruskalSorted for a version that relies on the edge list being already
// sorted.
func (l WeightedEdgeList) Kruskal() (f FromList, labels []int, dist float64) {
	sort.Sort(l)
	return l.KruskalSorted()
}

// KruskalSorted implements Kruskal's algorithm for constructing a minimum
// spanning tree on an undirected graph.
//
// When called, the edge list of the reciever must be already sorted by weight.
// See Kruskal for a version that accepts an unsorted edge list.
//
// The tree is returned as a FromList, with leaf nodes indicated, but without
// path lengths or MaxPathLen set.  Also returned is a parallel list of labels
// and a total distance for the returned tree.
func (l WeightedEdgeList) KruskalSorted() (f FromList, labels []int, dist float64) {
	ds := newDisjointSet(l.Order)
	// also initialize FromList as isolated nodes
	paths := make([]PathEnd, l.Order)
	for i := range paths {
		paths[i].From = -1
	}
	f.Paths = paths
	OneBits(&f.Leaves, l.Order)
	labels = make([]int, l.Order)
	// Kruskal for sorted edge list:
	for _, e := range l.Edges {
		x := e.N1
		y := e.N2
		xr := ds.find(x)
		yr := ds.find(y)
		if xr == yr {
			continue
		}
		switch xe, ye := &ds.set[xr], &ds.set[yr]; {
		case xe.rank < ye.rank:
			xe.from = yr
			x, y = y, x // swap so y is the "smaller" one
		case xe.rank == ye.rank:
			xe.rank++
			fallthrough
		default:
			ye.from = xr
		}
		// add arc so y comes from x
		paths[y].From = x
		labels[y] = e.Label
		f.Leaves.SetBit(&f.Leaves, x, 0)
		dist += l.WeightFunc(e.Label)
	}
	return
}

// Prim implements the Jarník-Prim-Dijkstra algorithm for constructing
// a minimum spanning tree on an undirected graph.
//
// Construct with NewPrim.
type Prim struct {
	Graph  LabeledAdjacencyList
	Weight WeightFunc
	Tree   FromList

	best []prNode // slice backs heap
}

// NewPrim constructs a new Prim object.  Argument g should represent an
// undirected graph.
func NewPrim(g LabeledAdjacencyList, w WeightFunc) *Prim {
	b := make([]prNode, len(g))
	for n := range b {
		b[n].nx = n
		b[n].fx = -1
	}
	return &Prim{
		Graph:  g,
		Weight: w,
		Tree:   NewFromList(len(g)),
		best:   b,
	}
}

type prNode struct {
	nx   int
	from FromHalf
	wt   float64 // p.Weight(from.Label)
	fx   int
}

type prHeap []*prNode

// Reset clears results of Span, allowing results to be recomputed.
//
// Reset is not meaningful following a change to the number of nodes in the
// graph.  To recompute following addition or deletion of nodes, simply
// abandon the Prim object and create a new one.
func (p *Prim) Reset() {
	p.Tree.reset()
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
	rp := p.Tree.Paths
	var frontier prHeap
	rp[start] = PathEnd{From: -1, Len: 1}
	b := p.best
	nDone := 1
	for a := start; ; {
		for _, nb := range p.Graph[a] {
			if rp[nb.To].Len > 0 {
				continue // already in MST, no action
			}
			switch bp := &b[nb.To]; {
			case bp.fx == -1: // new node for frontier
				bp.from = FromHalf{From: a, Label: nb.Label}
				bp.wt = p.Weight(nb.Label)
				heap.Push(&frontier, bp)
			case p.Weight(nb.Label) < bp.wt: // better arc
				bp.from = FromHalf{From: a, Label: nb.Label}
				bp.wt = p.Weight(nb.Label)
				heap.Fix(&frontier, bp.fx)
			}
		}
		if len(frontier) == 0 {
			break // done
		}
		bp := heap.Pop(&frontier).(*prNode)
		a = bp.nx
		rp[a].Len = rp[bp.from.From].Len + 1
		rp[a].From = bp.from.From
		//		p.Wt[a] = p.Weight(bp.from.Label)
		nDone++
	}
	return nDone
}

func (h prHeap) Len() int           { return len(h) }
func (h prHeap) Less(i, j int) bool { return h[i].wt < h[j].wt }
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
