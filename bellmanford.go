// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

type BellmanFord struct {
	Graph  WeightedAdjacencyList
	Result *WeightedFromTree
}

func NewBellmanFord(g WeightedAdjacencyList) *BellmanFord {
	return &BellmanFord{g, NewWeightedFromTree(len(g))}
}

func (b *BellmanFord) Run(start int) (ok bool) {
	b.Result.Reset()
	rp := b.Result.Paths
	for n := range rp {
		rp[n] = WeightedPathEnd{
			Dist: b.Result.NoPath,
			From: HalfFrom{-1, b.Result.NoPath},
		}
	}
	rp[start].Dist = 0
	rp[start].Len = 1
	for _ = range b.Graph[1:] {
		imp := false
		for from, nbs := range b.Graph {
			d1 := rp[from].Dist
			for _, nb := range nbs {
				d2 := d1 + nb.ArcWeight
				to := &rp[nb.To]
				if to.Len == 0 || d2 < to.Dist {
					*to = WeightedPathEnd{
						Dist: d2,
						From: HalfFrom{from, nb.ArcWeight},
						Len:  rp[from].Len + 1,
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
