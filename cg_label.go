// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"math/big"

	"github.com/willf/bitset"
)

// cg_adj.go is code generated from cg_label.go by directive in graph.go.
// Editing cg_label.go is okay.
// DO NOT EDIT cg_adj.go.

// ArcSize returns the number of arcs in g.
//
// Note that for an undirected graph without loops, the number of undirected
// edges -- the traditional meaning of graph size -- will be ArcSize()/2.
// On the other hand, if g is an undirected graph that has or may have loops,
// g.ArcSize()/2 is not a meaningful quantity.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) ArcSize() int {
	m := 0
	for _, to := range g {
		m += len(to)
	}
	return m
}

// Balanced returns true if for every node in g, in-degree equals out-degree.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) Balanced() bool {
	for n, in := range g.InDegree() {
		if in != len(g[n]) {
			return false
		}
	}
	return true
}

// Bipartite determines if a connected component of an undirected graph
// is bipartite.
//
// Argument n can be any representative node of the component.
//
// If the component is bipartite, Bipartite returns true and a two-coloring
// of the component.  Each color set is returned as a bitmap.  If the component
// is not bipartite, Bipartite returns false and a representative odd cycle.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) Bipartite(n NI) (b bool, c1, c2 *big.Int, oc []NI) {
	c1 = &big.Int{}
	c2 = &big.Int{}
	b = true
	var open bool
	var df func(n NI, c1, c2 *big.Int)
	df = func(n NI, c1, c2 *big.Int) {
		c1.SetBit(c1, int(n), 1)
		for _, nb := range g[n] {
			if c1.Bit(int(nb.To)) == 1 {
				b = false
				oc = []NI{nb.To, n}
				open = true
				return
			}
			if c2.Bit(int(nb.To)) == 1 {
				continue
			}
			df(nb.To, c2, c1)
			if b {
				continue
			}
			switch {
			case !open:
			case n == oc[0]:
				open = false
			default:
				oc = append(oc, n)
			}
			return
		}
	}
	df(n, c1, c2)
	if b {
		return b, c1, c2, nil
	}
	return b, nil, nil, oc
}

// BoundsOk validates that all arcs in g stay within the slice bounds of g.
//
// BoundsOk returns true when no arcs point outside the bounds of g.
// Otherwise it returns false and an example arc that points outside of g.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) BoundsOk() (ok bool, fr NI, to Half) {
	for fr, to := range g {
		for _, to := range to {
			if to.To < 0 || to.To >= NI(len(g)) {
				return false, NI(fr), to
			}
		}
	}
	return true, -1, to
}

// BronKerbosch1 finds maximal cliques in an undirected graph.
//
// The graph must not contain parallel edges or loops.
//
// See https://en.wikipedia.org/wiki/Clique_(graph_theory) and
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm for background.
//
// This method implements the BronKerbosch1 algorithm of WP; that is,
// the original algorithm without improvements.
//
// The method sends all maximal cliques in g on the returned channel, then
// closes the channel.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also more sophisticated variants BronKerbosch2 and BronKerbosch3.
func (g LabeledAdjacencyList) BronKerbosch1() chan []int {
	ch := make(chan []int)
	go func() {
		var f func(R, P, X *bitset.BitSet)
		f = func(R, P, X *bitset.BitSet) {
			switch {
			case P.Any():
				r2 := bitset.New(uint(len(g)))
				p2 := bitset.New(uint(len(g)))
				x2 := bitset.New(uint(len(g)))
				for n, ok := P.NextSet(0); ok; n, ok = P.NextSet(n + 1) {
					R.Copy(r2)
					r2.Set(n)
					p2.ClearAll()
					x2.ClearAll()
					for _, to := range g[n] {
						if P.Test(uint(to.To)) {
							p2.Set(uint(to.To))
						}
						if X.Test(uint(to.To)) {
							x2.Set(uint(to.To))
						}
					}
					f(r2, p2, x2)
					P.SetTo(n, false)
					X.Set(n)
				}
			case X.None():
				var n uint
				n--
				c := make([]int, R.Count())
				for i := range c {
					n, _ = R.NextSet(n + 1)
					c[i] = int(n)
				}
				ch <- c
			}
		}
		R := bitset.New(uint(len(g)))
		P := bitset.New(uint(len(g))).Complement()
		X := bitset.New(uint(len(g)))
		f(R, P, X)
		close(ch)
	}()
	return ch
}

// BKPivotMinP is a strategy for BronKerbosch methods.
//
// To use it, take the method value (see golang.org/ref/spec#Method_values)
// and pass it as the argument to BronKerbosch2 or 3.
//
// The strategy is to simply pick the first node in P.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) BKPivotMinP(P, X *bitset.BitSet) int {
	n, _ := P.NextSet(0)
	return int(n)
}

// BKPivotMaxDegree is a strategy for BronKerbosch methods.
//
// To use it, take the method value (see golang.org/ref/spec#Method_values)
// and pass it as the argument to BronKerbosch2 or 3.
//
// The strategy is to pick the node from P or X with the maximum degree
// (number of edges) in g.  Note this is a shortcut from evaluating degrees
// in P.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) BKPivotMaxDegree(P, X *bitset.BitSet) int {
	// choose pivot u as highest degree node from P or X
	n, ok := P.NextSet(0)
	u := n
	maxDeg := len(g[u])
	for { // scan P
		n, ok = P.NextSet(n + 1)
		if !ok {
			break
		}
		if d := len(g[n]); d > maxDeg {
			u = n
			maxDeg = d
		}
	}
	// scan X
	for n, ok = X.NextSet(0); ok; n, ok = X.NextSet(n + 1) {
		if d := len(g[n]); d > maxDeg {
			u = n
			maxDeg = d
		}
	}
	return int(u)
}

// BronKerbosch2 finds maximal cliques in an undirected graph.
//
// The graph must not contain parallel edges or loops.
//
// See https://en.wikipedia.org/wiki/Clique_(graph_theory) and
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm for background.
//
// This method implements the BronKerbosch2 algorithm of WP; that is,
// the original algorithm plus pivoting.
//
// The argument is a pivot function that must return a node of P or X.
// P is guaranteed to contain at least one node.  X is not.
// For example see BKPivotMaxDegree.
//
// The method sends all maximal cliques in g on the returned channel, then
// closes the channel.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also simpler variant BronKerbosch1 and more sophisticated variant
// BronKerbosch3.
func (g LabeledAdjacencyList) BronKerbosch2(pivot func(P, X *bitset.BitSet) int) chan []int {
	ch := make(chan []int)
	go func() {
		var f func(R, P, X *bitset.BitSet)
		f = func(R, P, X *bitset.BitSet) {
			switch {
			case P.Any():
				r2 := bitset.New(uint(len(g)))
				p2 := bitset.New(uint(len(g)))
				x2 := bitset.New(uint(len(g)))
				// compute P \ N(u).  next 5 lines are only difference from BK1
				pnu := P.Clone()
				for _, to := range g[pivot(P, X)] {
					pnu.SetTo(uint(to.To), false)
				}
				for n, ok := pnu.NextSet(0); ok; n, ok = pnu.NextSet(n + 1) {
					// remaining code like BK1
					R.Copy(r2)
					r2.Set(n)
					p2.ClearAll()
					x2.ClearAll()
					for _, to := range g[n] {
						if P.Test(uint(to.To)) {
							p2.Set(uint(to.To))
						}
						if X.Test(uint(to.To)) {
							x2.Set(uint(to.To))
						}
					}
					f(r2, p2, x2)
					P.SetTo(n, false)
					X.Set(n)
				}
			case X.None():
				var n uint
				n--
				c := make([]int, R.Count())
				for i := range c {
					n, _ = R.NextSet(n + 1)
					c[i] = int(n)
				}
				ch <- c
			}
		}
		R := bitset.New(uint(len(g)))
		P := bitset.New(uint(len(g))).Complement()
		X := bitset.New(uint(len(g)))
		f(R, P, X)
		close(ch)
	}()
	return ch
}

// BronKerbosch3 finds maximal cliques in an undirected graph.
//
// The graph must not contain parallel edges or loops.
//
// See https://en.wikipedia.org/wiki/Clique_(graph_theory) and
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm for background.
//
// This method implements the BronKerbosch3 algorithm of WP; that is,
// the original algorithm with pivoting and degeneracy ordering.
//
// The argument is a pivot function that must return a node of P or X.
// P is guaranteed to contain at least one node.  X is not.
// For example see BKPivotMaxDegree.
//
// The method sends all maximal cliques in g on the returned channel, then
// closes the channel.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also simpler variants BronKerbosch1 and BronKerbosch2.
func (g LabeledAdjacencyList) BronKerbosch3(pivot func(P, X *bitset.BitSet) int) chan []int {
	ch := make(chan []int)
	go func() {
		var f func(R, P, X *bitset.BitSet)
		f = func(R, P, X *bitset.BitSet) {
			switch {
			case P.Any():
				r2 := bitset.New(uint(len(g)))
				p2 := bitset.New(uint(len(g)))
				x2 := bitset.New(uint(len(g)))
				// compute P \ N(u).  next 5 lines are only difference from BK1
				pnu := P.Clone()
				for _, to := range g[pivot(P, X)] {
					pnu.SetTo(uint(to.To), false)
				}
				for n, ok := pnu.NextSet(0); ok; n, ok = pnu.NextSet(n + 1) {
					// remaining code like BK1
					R.Copy(r2)
					r2.Set(n)
					p2.ClearAll()
					x2.ClearAll()
					for _, to := range g[n] {
						if P.Test(uint(to.To)) {
							p2.Set(uint(to.To))
						}
						if X.Test(uint(to.To)) {
							x2.Set(uint(to.To))
						}
					}
					f(r2, p2, x2)
					P.SetTo(n, false)
					X.Set(n)
				}
			case X.None():
				var n uint
				n--
				c := make([]int, R.Count())
				for i := range c {
					n, _ = R.NextSet(n + 1)
					c[i] = int(n)
				}
				ch <- c
			}
		}

		R := bitset.New(uint(len(g)))
		P := bitset.New(uint(len(g))).Complement()
		X := bitset.New(uint(len(g)))
		// code above same as BK2
		// code below new to BK3
		_, ord, _ := g.Degeneracy()
		p2 := bitset.New(uint(len(g)))
		x2 := bitset.New(uint(len(g)))
		for _, n := range ord {
			R.Set(uint(n))
			p2.ClearAll()
			x2.ClearAll()
			for _, to := range g[n] {
				if P.Test(uint(to.To)) {
					p2.Set(uint(to.To))
				}
				if X.Test(uint(to.To)) {
					x2.Set(uint(to.To))
				}
			}
			f(R, p2, x2)
			R.SetTo(uint(n), false)
			P.SetTo(uint(n), false)
			X.Set(uint(n))
		}
		close(ch)
	}()
	return ch
}

// ConnectedComponentBits, for undirected graphs, returns a function that
// iterates over connected components of g, returning a member bitmap for each.
//
// Each call of the returned function returns the order (number of nodes)
// and bits of a connected component.  The returned function returns zeros
// after returning all connected components.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps, which has lighter weight return values.
func (g LabeledAdjacencyList) ConnectedComponentBits() func() (order int, bits big.Int) {
	var vg big.Int  // nodes visited in graph
	var vc *big.Int // nodes visited in current component
	var nc int
	var df func(NI)
	df = func(n NI) {
		vg.SetBit(&vg, int(n), 1)
		vc.SetBit(vc, int(n), 1)
		nc++
		for _, nb := range g[n] {
			if vg.Bit(int(nb.To)) == 0 {
				df(nb.To)
			}
		}
		return
	}
	var n NI
	return func() (o int, bits big.Int) {
		for ; n < NI(len(g)); n++ {
			if vg.Bit(int(n)) == 0 {
				vc = &bits
				nc = 0
				df(n)
				return nc, bits
			}
		}
		return
	}
}

// ConnectedComponentLists, for undirected graphs, returns a function that
// iterates over connected components of g, returning the member list of each.
//
// Each call of the returned function returns a node list of a
// connected component.  The returned function returns nil after returning
// all connected components.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps, which has lighter weight return values.
func (g LabeledAdjacencyList) ConnectedComponentLists() func() []NI {
	var vg big.Int // nodes visited in graph
	var m []NI     // members of current component
	var df func(NI)
	df = func(n NI) {
		vg.SetBit(&vg, int(n), 1)
		m = append(m, n)
		for _, nb := range g[n] {
			if vg.Bit(int(nb.To)) == 0 {
				df(nb.To)
			}
		}
		return
	}
	var n NI
	return func() []NI {
		for ; n < NI(len(g)); n++ {
			if vg.Bit(int(n)) == 0 {
				m = nil
				df(n)
				return m
			}
		}
		return nil
	}
}

// ConnectedComponentReps, for undirected graphs, returns a representative
// node from each connected component of g.
//
// Returned is a slice with a single representative node from each connected
// component and also a parallel slice with the order, or number of nodes,
// in the corresponding component.
//
// This is fairly minimal information describing connected components.
// From a representative node, other nodes in the component can be reached
// by depth first traversal for example.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentBits and ConnectedComponentLists which can
// collect component members in a single traversal.
func (g LabeledAdjacencyList) ConnectedComponentReps() (reps []NI, orders []int) {
	var c big.Int
	var o int
	var df func(NI)
	df = func(n NI) {
		c.SetBit(&c, int(n), 1)
		o++
		for _, nb := range g[n] {
			if c.Bit(int(nb.To)) == 0 {
				df(nb.To)
			}
		}
		return
	}
	for n := range g {
		if c.Bit(n) == 0 {
			reps = append(reps, NI(n))
			o = 0
			df(NI(n))
			orders = append(orders, o)
		}
	}
	return
}

// Copy makes a copy of g, copying the underlying slices.
// Copy also computes the arc size m, the number of arcs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) Copy() (c LabeledAdjacencyList, m int) {
	c = make(LabeledAdjacencyList, len(g))
	for n, to := range g {
		c[n] = append([]Half{}, to...)
		m += len(to)
	}
	return
}

// Cyclic, for directed graphs, determines if g contains cycles.
//
// Cyclic returns true if g contains at least one cycle.
// Cyclic returns false if g is acyclic.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) Cyclic() (cyclic bool, fr NI, to Half) {
	fr, to.To = -1, -1
	var temp, perm big.Int
	var df func(NI)
	df = func(n NI) {
		switch {
		case temp.Bit(int(n)) == 1:
			cyclic = true
			return
		case perm.Bit(int(n)) == 1:
			return
		}
		temp.SetBit(&temp, int(n), 1)
		for _, nb := range g[n] {
			df(nb.To)
			if cyclic {
				if fr < 0 {
					fr, to = n, nb
				}
				return
			}
		}
		temp.SetBit(&temp, int(n), 0)
		perm.SetBit(&perm, int(n), 1)
	}
	for n := range g {
		if perm.Bit(n) == 1 {
			continue
		}
		if df(NI(n)); cyclic { // short circuit as soon as a cycle is found
			break
		}
	}
	return
}

// Degeneracy computes k-degeneracy, vertex ordering and k-cores.
//
// See Wikipedia https://en.wikipedia.org/wiki/Degeneracy_(graph_theory)
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) Degeneracy() (k int, ordering []NI, cores []int) {
	// WP algorithm
	ordering = make([]NI, len(g))
	var L big.Int
	d := make([]int, len(g))
	var D [][]NI
	for v, nb := range g {
		dv := len(nb)
		d[v] = dv
		for len(D) <= dv {
			D = append(D, nil)
		}
		D[dv] = append(D[dv], NI(v))
	}
	for ox := range g {
		// find a non-empty D
		i := 0
		for len(D[i]) == 0 {
			i++
		}
		// k is max(i, k)
		if i > k {
			for len(cores) <= i {
				cores = append(cores, 0)
			}
			cores[k] = ox
			k = i
		}
		// select from D[i]
		Di := D[i]
		last := len(Di) - 1
		v := Di[last]
		// Add v to ordering, remove from Di
		ordering[ox] = v
		L.SetBit(&L, int(v), 1)
		D[i] = Di[:last]
		// move neighbors
		for _, nb := range g[v] {
			if L.Bit(int(nb.To)) == 1 {
				continue
			}
			dn := d[nb.To] // old number of neighbors of nb
			Ddn := D[dn]   // nb is in this list
			// remove it from the list
			for wx, w := range Ddn {
				if w == nb.To {
					last := len(Ddn) - 1
					Ddn[wx], Ddn[last] = Ddn[last], Ddn[wx]
					D[dn] = Ddn[:last]
				}
			}
			dn-- // new number of neighbors
			d[nb.To] = dn
			// re--add it to it's new list
			D[dn] = append(D[dn], nb.To)
		}
	}
	cores[k] = len(ordering)
	return
}

// DepthFirst traverses a graph depth first.
//
// As it traverses it calls visitor function v for each node.  If v returns
// false at any point, the traversal is terminated immediately and DepthFirst
// returns false.  Otherwise DepthFirst returns true.
//
// DepthFirst uses argument bm is used as a bitmap to guide the traversal.
// For a complete traversal, bm should be 0 initially.  During the
// traversal, bits are set corresponding to each node visited.
// The bit is set before calling the visitor function.
//
// Argument bm can be nil if you have no need for it.
// In this case a bitmap is created internally for one-time use.
//
// Alternatively v can be nil.  In this case traversal still procedes and
// updates the bitmap, which can be a useful result.
// DepthFirst always returns true in this case.
//
// It makes no sense for both bm and v to be nil.  In this case DepthFirst
// returns false immediately.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) DepthFirst(start NI, bm *big.Int, v Visitor) (ok bool) {
	if bm == nil {
		if v == nil {
			return false
		}
		bm = new(big.Int)
	}
	ok = true
	var df func(n NI)
	df = func(n NI) {
		if bm.Bit(int(n)) == 1 {
			return
		}
		bm.SetBit(bm, int(n), 1)
		if v != nil && !v(n) {
			ok = false
			return
		}
		for _, nb := range g[n] {
			df(nb.To)
		}
	}
	df(start)
	return
}

// FromList transposes a graph into a FromList, typically to encode a tree.
//
// Results may not be meaningful for non-trees.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) FromList() FromList {
	// init paths
	paths := make([]PathEnd, len(g))
	for i := range paths {
		paths[i].From = -1
	}
	// init leaves
	var leaves big.Int
	OneBits(&leaves, len(g))
	// iterate over arcs, setting from pointers and and marking non-leaves.
	for fr, to := range g {
		for _, to := range to {
			paths[to.To].From = NI(fr)
			leaves.SetBit(&leaves, fr, 0)
		}
	}
	// f to set path lengths
	var leng func(NI) int
	leng = func(n NI) int {
		if l := paths[n].Len; l > 0 {
			return l
		}
		fr := paths[n].From
		if fr < 0 {
			paths[n].Len = 1
			return 1
		}
		l := 1 + leng(fr)
		paths[n].Len = l
		return l
	}
	// for each leaf, trace path to set path length, accumulate max
	maxLen := 0
	for i := range paths {
		if leaves.Bit(i) == 1 {
			if l := leng(NI(i)); l > maxLen {
				maxLen = l
			}
		}
	}
	return FromList{paths, leaves, maxLen}
}

// HasLoop identifies if a graph contains a loop, an arc that leads from a
// a node back to the same node.
//
// If the graph has a loop, the result is an example node that has a loop.
//
// If there are no loops, the method returns -1.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) HasLoop() NI {
	for fr, to := range g {
		for _, to := range to {
			if NI(fr) == to.To {
				return to.To
			}
		}
	}
	return -1
}

// HasParallelMap identifies if a graph contains parallel arcs, multiple arcs
// that lead from a node to the same node.
//
// If the graph has parallel arcs, the results fr and to represent an example
// where there are parallel arcs from node fr to node to.
//
// If there are no parallel arcs, the method returns -1 -1.
//
// Multiple loops on a node count as parallel arcs.
//
// "Map" in the method name indicates that a Go map is used to detect parallel
// arcs.  Compared to method HasParallelSort, this gives better asymtotic
// performance for large dense graphs but may have increased overhead for
// small or sparse graphs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) HasParallelMap() (fr, to NI) {
	for n, to := range g {
		if len(to) == 0 {
			continue
		}
		m := map[NI]struct{}{}
		for _, to := range to {
			if _, ok := m[to.To]; ok {
				return NI(n), to.To
			}
			m[to.To] = struct{}{}
		}
	}
	return -1, -1
}

// InDegree computes the in-degree of each node in g
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) InDegree() []int {
	ind := make([]int, len(g))
	for _, nbs := range g {
		for _, nb := range nbs {
			ind[nb.To]++
		}
	}
	return ind
}

// IsSimple checks for loops and parallel arcs.
//
// A graph is "simple" if it has no loops or parallel arcs.
//
// IsSimple returns true, -1 for simple graphs.  If a loop or parallel arc is
// found, simple returns false and a node that represents a counterexample
// to the graph being simple.
//
// See also separate methods HasLoop and HasParallel.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) IsSimple() (ok bool, n NI) {
	if n = g.HasLoop(); n >= 0 {
		return
	}
	if n, _ = g.HasParallelSort(); n >= 0 {
		return
	}
	return true, -1
}

// IsTreeDirected identifies trees in directed graphs.
//
// IsTreeDirected returns true if the subgraph reachable from
// root is a tree.  It does not validate that the entire graph is a tree.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) IsTreeDirected(root NI) bool {
	var v big.Int
	var df func(NI) bool
	df = func(n NI) bool {
		if v.Bit(int(n)) == 1 {
			return false
		}
		v.SetBit(&v, int(n), 1)
		for _, to := range g[n] {
			if !df(to.To) {
				return false
			}
		}
		return true
	}
	return df(root)
}

// IsTreeUndirected identifies trees in undirected graphs.
//
// IsTreeUndirected returns true if the connected component
// containing argument root is a tree.  It does not validate
// that the entire graph is a tree.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) IsTreeUndirected(root NI) bool {
	var v big.Int
	var df func(NI, NI) bool
	df = func(fr, n NI) bool {
		if v.Bit(int(n)) == 1 {
			return false
		}
		v.SetBit(&v, int(n), 1)
		for _, to := range g[n] {
			if to.To != fr && !df(n, to.To) {
				return false
			}
		}
		return true
	}
	v.SetBit(&v, int(root), 1)
	for _, to := range g[root] {
		if !df(root, to.To) {
			return false
		}
	}
	return true
}

/*
MaxmimalClique finds a maximal clique containing the node n.

Not sure this is good for anything.  It produces a single maximal clique
but there can be multiple maximal cliques containing a given node.
This algorithm just returns one of them, not even necessarily the
largest one.

func (g LabeledAdjacencyList) MaximalClique(n int) []int {
	c := []int{n}
	var m bitset.BitSet
	m.Set(uint(n))
	for fr, to := range g {
		if fr == n {
			continue
		}
		if len(to) < len(c) {
			continue
		}
		f := 0
		for _, to := range to {
			if m.Test(uint(to.To)) {
				f++
				if f == len(c) {
					c = append(c, to.To)
					m.Set(uint(to.To))
					break
				}
			}
		}
	}
	return c
}
*/

// Tarjan identifies strongly connected components in a directed graph using
// Tarjan's algorithm.
//
// Returned is a list of components, each component is a list of nodes.
// A property of the algorithm is that components are ordered in reverse
// topological order of the condensation.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also TarjanForward and TarjanCondensation.
func (g LabeledAdjacencyList) Tarjan() (scc [][]NI) {
	// See "Depth-first search and linear graph algorithms", Robert Tarjan,
	// SIAM J. Comput. Vol. 1, No. 2, June 1972.
	//
	// Implementation here from Wikipedia pseudocode,
	// http://en.wikipedia.org/w/index.php?title=Tarjan%27s_strongly_connected_components_algorithm&direction=prev&oldid=647184742
	//
	// There are equivalent labeled and unlabeled versions of this method.
	var indexed, stacked big.Int
	index := make([]int, len(g))
	lowlink := make([]int, len(g))
	x := 0
	var S []NI
	var sc func(NI)
	sc = func(n NI) {
		index[n] = x
		indexed.SetBit(&indexed, int(n), 1)
		lowlink[n] = x
		x++
		S = append(S, n)
		stacked.SetBit(&stacked, int(n), 1)
		for _, nb := range g[n] {
			if indexed.Bit(int(nb.To)) == 0 {
				sc(nb.To)
				if lowlink[nb.To] < lowlink[n] {
					lowlink[n] = lowlink[nb.To]
				}
			} else if stacked.Bit(int(nb.To)) == 1 {
				if index[nb.To] < lowlink[n] {
					lowlink[n] = index[nb.To]
				}
			}
		}
		if lowlink[n] == index[n] {
			var c []NI
			for {
				last := len(S) - 1
				w := S[last]
				S = S[:last]
				stacked.SetBit(&stacked, int(w), 0)
				c = append(c, w)
				if w == n {
					scc = append(scc, c)
					break
				}
			}
		}
	}
	for n := range g {
		if indexed.Bit(n) == 0 {
			sc(NI(n))
		}
	}
	return
}

// TarjanForward returns strongly connected components.
//
// It returns components in the reverse order of Tarjan, for situations
// where a forward topological ordering is easier.
func (g LabeledAdjacencyList) TarjanForward() (scc [][]NI) {
	scc = g.Tarjan()
	last := len(scc) - 1
	for i, ci := range scc[:len(scc)/2] {
		scc[i], scc[last-i] = scc[last-i], ci
	}
	return
}

// TarjanCondensation returns strongly connected components and their
// condensation graph.
//
// Components are ordered in a forward topological ordering.
func (g LabeledAdjacencyList) TarjanCondensation() (scc [][]NI, cd AdjacencyList) {
	scc = g.TarjanForward()
	cd = make(AdjacencyList, len(scc)) // return value
	cond := make([]NI, len(g))         // mapping from g node to cd node
	for cn := len(scc) - 1; cn >= 0; cn-- {
		c := scc[cn]
		for _, n := range c {
			cond[n] = NI(cn) // map g node to cd node
		}
		var tos []NI  // list of 'to' nodes
		var m big.Int // tos map
		m.SetBit(&m, cn, 1)
		for _, n := range c {
			for _, to := range g[n] {
				if ct := cond[to.To]; m.Bit(int(ct)) == 0 {
					m.SetBit(&m, int(ct), 1)
					tos = append(tos, ct)
				}
			}
		}
		cd[cn] = tos
	}
	return
}

// Topological, for directed acyclic graphs, computes a topological sort of g.
//
// For an acyclic graph, return value ordering is a permutation of node numbers
// in topologically sorted order and cycle will be nil.  If the graph is found
// to be cyclic, ordering will be nil and cycle will be the path of a found
// cycle.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) Topological() (ordering, cycle []NI) {
	ordering = make([]NI, len(g))
	i := len(ordering)
	var temp, perm big.Int
	var cycleFound bool
	var cycleStart NI
	var df func(NI)
	df = func(n NI) {
		switch {
		case temp.Bit(int(n)) == 1:
			cycleFound = true
			cycleStart = n
			return
		case perm.Bit(int(n)) == 1:
			return
		}
		temp.SetBit(&temp, int(n), 1)
		for _, nb := range g[n] {
			df(nb.To)
			if cycleFound {
				if cycleStart >= 0 {
					// a little hack: orderng won't be needed so repurpose the
					// slice as cycle.  this is read out in reverse order
					// as the recursion unwinds.
					x := len(ordering) - 1 - len(cycle)
					ordering[x] = n
					cycle = ordering[x:]
					if n == cycleStart {
						cycleStart = -1
					}
				}
				return
			}
		}
		temp.SetBit(&temp, int(n), 0)
		perm.SetBit(&perm, int(n), 1)
		i--
		ordering[i] = n
	}
	for n := range g {
		if perm.Bit(n) == 1 {
			continue
		}
		df(NI(n))
		if cycleFound {
			return nil, cycle
		}
	}
	return ordering, nil
}

// TopologicalKahn, for directed acyclic graphs, computes a topological sort of g.
//
// For an acyclic graph, return value ordering is a permutation of node numbers
// in topologically sorted order and cycle will be nil.  If the graph is found
// to be cyclic, ordering will be nil and cycle will be the path of a found
// cycle.
//
// This function is based on the algorithm by Arthur Kahn and requires the
// transpose of g be passed as the argument.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) TopologicalKahn(tr AdjacencyList) (ordering, cycle []NI) {
	// code follows Wikipedia pseudocode.
	var L, S []NI
	// rem for "remaining edges," this function makes a local copy of the
	// in-degrees and consumes that instead of consuming an input.
	rem := make([]int, len(g))
	for n, fr := range tr {
		if len(fr) == 0 {
			// accumulate "set of all nodes with no incoming edges"
			S = append(S, NI(n))
		} else {
			// initialize rem from in-degree
			rem[n] = len(fr)
		}
	}
	for len(S) > 0 {
		last := len(S) - 1 // "remove a node n from S"
		n := S[last]
		S = S[:last]
		L = append(L, n) // "add n to tail of L"
		for _, m := range g[n] {
			// WP pseudo code reads "for each node m..." but it means for each
			// node m *remaining in the graph.*  We consume rem rather than
			// the graph, so "remaining in the graph" for us means rem[m] > 0.
			if rem[m.To] > 0 {
				rem[m.To]--         // "remove edge from the graph"
				if rem[m.To] == 0 { // if "m has no other incoming edges"
					S = append(S, m.To) // "insert m into S"
				}
			}
		}
	}
	// "If graph has edges," for us means a value in rem is > 0.
	for c, in := range rem {
		if in > 0 {
			// recover cyclic nodes
			for _, nb := range g[c] {
				if rem[nb.To] > 0 {
					cycle = append(cycle, NI(c))
					break
				}
			}
		}
	}
	if len(cycle) > 0 {
		return nil, cycle
	}
	return L, nil
}
