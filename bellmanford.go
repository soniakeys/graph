// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"math"
)

type BellmanFord struct {
	Graph  AdjacencyList // input graph
	Result []BellmanFordResult
}

type BellmanFordResult struct {
	Dist float64
	From FromHalf
}

func NewBellmanFord(g AdjacencyList) *BellmanFord {
	return &BellmanFord{Graph: g, Result: make([]BellmanFordResult, len(g))}
}

func (b *BellmanFord) Scan(start int) (ok bool) {
	r := b.Result
	for n := range r {
		r[n] = BellmanFordResult{math.Inf(-1), FromHalf{-1, math.Inf(1)}}
	}
	for _ = range b.Graph[1:] {
		imp := false
		for from, nbs := range b.Graph {
			d1 := r[from].Dist
			for _, nb := range nbs {
				if d2 := d1 + nb.ArcWeight; d2 < r[nb.To].Dist {
					r[nb.To] = BellmanFordResult{
						d2,
						FromHalf{from, nb.ArcWeight},
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
		d1 := r[from].Dist
		for _, nb := range nbs {
			if d1+nb.ArcWeight < r[nb.To].Dist {
				return false // negative cycle
			}
		}
	}
	return true
}
