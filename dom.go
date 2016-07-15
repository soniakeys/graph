// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

type Dominators struct {
	Immediate []NI
	from      Directed
}

func (d Dominators) Set(n NI) []NI {
	im := d.Immediate
	if im[n] < 0 {
		return nil
	}
	for s := []NI{n}; ; {
		if p := im[n]; p < 0 || p == n {
			return s
		} else {
			s = append(s, p)
			n = p
		}
	}
}

func (d Dominators) Frontier() []map[NI]struct{} {
	im := d.Immediate
	f := make([]map[NI]struct{}, len(im))
	for i := range f {
		f[i] = map[NI]struct{}{}
	}
	for b, fr := range d.from.AdjacencyList {
		if len(fr) < 2 {
			continue
		}
		imb := im[b]
		for _, p := range fr {
			for runner := p; runner != imb; runner = im[runner] {
				f[runner][NI(b)] = struct{}{}
			}
		}
	}
	return f
}

// returns immediate dominator for each node.  (a from-list)
func (g Directed) Dominators(start NI) Dominators {
	a := g.AdjacencyList
	post := make([]NI, len(a))
	l := len(a)
	a.DepthFirst(start, nil, func(n NI) bool {
		l--
		post[l] = n
		return true
	})
	tr, _ := g.Transpose()
	return g.Doms(tr, post[l:])
}

func (g Directed) Doms(tr Directed, post []NI) Dominators {
	a := g.AdjacencyList
	dom := make([]NI, len(a))
	pi := make([]int, len(a))
	for i, n := range post {
		pi[n] = i
	}
	intersect := func(b1, b2 NI) NI {
		for b1 != b2 {
			for pi[b1] < pi[b2] {
				b1 = dom[b1]
			}
			for pi[b2] < pi[b1] {
				b2 = dom[b2]
			}
		}
		return b1
	}
	for n := range dom {
		dom[n] = -1
	}
	start := post[len(post)-1]
	dom[start] = start
	for changed := false; ; changed = false {
		for i := len(post) - 2; i >= 0; i-- {
			b := post[i]
			var im, fp NI
			fr := tr.AdjacencyList[b]
			var j int
			for j, fp = range fr {
				if dom[fp] >= 0 {
					im = fp
					break
				}
			}
			for _, p := range fr[j:] {
				if dom[p] >= 0 {
					im = intersect(im, p)
				}
			}
			if dom[b] != im {
				dom[b] = im
				changed = true
			}
		}
		if !changed {
			return Dominators{dom, tr}
		}
	}
}

func (g Directed) PostDominators(end NI) Dominators {
	tr, _ := g.Transpose()
	a := tr.AdjacencyList
	post := make([]NI, len(a))
	l := len(a)
	a.DepthFirst(end, nil, func(n NI) bool {
		l--
		post[l] = n
		return true
	})
	return tr.Doms(g, post[l:])
}
