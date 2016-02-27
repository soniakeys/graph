// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"math/big"

	"github.com/willf/bitset"
)

// cg_undir_al.go is code generated from cg_undir.go by directive in graph.go.
// Editing cg_undir.go is okay.
// DO NOT EDIT cg_undir_al.go.

// Bipartite determines if a connected component of an undirected graph
// is bipartite, a component where nodes can be partitioned into two sets
// such that every edge in the component goes from one set to the other.
//
// Argument n can be any representative node of the component.
//
// If the component is bipartite, Bipartite returns true and a two-coloring
// of the component.  Each color set is returned as a bitmap.  If the component
// is not bipartite, Bipartite returns false and a representative odd cycle.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g Undirected) Bipartite(n NI) (b bool, c1, c2 *big.Int, oc []NI) {
	c1 = &big.Int{}
	c2 = &big.Int{}
	b = true
	var open bool
	var df func(n NI, c1, c2 *big.Int)
	df = func(n NI, c1, c2 *big.Int) {
		c1.SetBit(c1, int(n), 1)
		for _, nb := range g.AdjacencyList[n] {
			if c1.Bit(int(nb)) == 1 {
				b = false
				oc = []NI{nb, n}
				open = true
				return
			}
			if c2.Bit(int(nb)) == 1 {
				continue
			}
			df(nb, c2, c1)
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
// The method calls the emit argument for each maximal clique in g, as long
// as emit returns true.  If emit returns false, BronKerbosch1 returns
// immediately.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also more sophisticated variants BronKerbosch2 and BronKerbosch3.
func (g Undirected) BronKerbosch1(emit func([]NI) bool) {
	var f func(R, P, X *bitset.BitSet) bool
	f = func(R, P, X *bitset.BitSet) bool {
		switch {
		case P.Any():
			r2 := bitset.New(uint(len(g.AdjacencyList)))
			p2 := bitset.New(uint(len(g.AdjacencyList)))
			x2 := bitset.New(uint(len(g.AdjacencyList)))
			for n, ok := P.NextSet(0); ok; n, ok = P.NextSet(n + 1) {
				R.Copy(r2)
				r2.Set(n)
				p2.ClearAll()
				x2.ClearAll()
				for _, to := range g.AdjacencyList[n] {
					if P.Test(uint(to)) {
						p2.Set(uint(to))
					}
					if X.Test(uint(to)) {
						x2.Set(uint(to))
					}
				}
				if !f(r2, p2, x2) {
					return false
				}
				P.SetTo(n, false)
				X.Set(n)
			}
		case X.None():
			var n uint
			n--
			c := make([]NI, R.Count())
			for i := range c {
				n, _ = R.NextSet(n + 1)
				c[i] = NI(n)
			}
			if !emit(c) {
				return false
			}
		}
		return true
	}
	R := bitset.New(uint(len(g.AdjacencyList)))
	P := bitset.New(uint(len(g.AdjacencyList))).Complement()
	X := bitset.New(uint(len(g.AdjacencyList)))
	f(R, P, X)
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
func (g Undirected) BKPivotMaxDegree(P, X *bitset.BitSet) int {
	// choose pivot u as highest degree node from P or X
	n, ok := P.NextSet(0)
	u := n
	maxDeg := len(g.AdjacencyList[u])
	for { // scan P
		n, ok = P.NextSet(n + 1)
		if !ok {
			break
		}
		if d := len(g.AdjacencyList[n]); d > maxDeg {
			u = n
			maxDeg = d
		}
	}
	// scan X
	for n, ok = X.NextSet(0); ok; n, ok = X.NextSet(n + 1) {
		if d := len(g.AdjacencyList[n]); d > maxDeg {
			u = n
			maxDeg = d
		}
	}
	return int(u)
}

// BKPivotMinP is a strategy for BronKerbosch methods.
//
// To use it, take the method value (see golang.org/ref/spec#Method_values)
// and pass it as the argument to BronKerbosch2 or 3.
//
// The strategy is to simply pick the first node in P.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g Undirected) BKPivotMinP(P, X *bitset.BitSet) int {
	n, _ := P.NextSet(0)
	return int(n)
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
// The method calls the emit argument for each maximal clique in g, as long
// as emit returns true.  If emit returns false, BronKerbosch1 returns
// immediately.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also simpler variant BronKerbosch1 and more sophisticated variant
// BronKerbosch3.
func (g Undirected) BronKerbosch2(pivot func(P, X *bitset.BitSet) int, emit func([]NI) bool) {
	var f func(R, P, X *bitset.BitSet) bool
	f = func(R, P, X *bitset.BitSet) bool {
		switch {
		case P.Any():
			r2 := bitset.New(uint(len(g.AdjacencyList)))
			p2 := bitset.New(uint(len(g.AdjacencyList)))
			x2 := bitset.New(uint(len(g.AdjacencyList)))
			// compute P \ N(u).  next 5 lines are only difference from BK1
			pnu := P.Clone()
			for _, to := range g.AdjacencyList[pivot(P, X)] {
				pnu.SetTo(uint(to), false)
			}
			for n, ok := pnu.NextSet(0); ok; n, ok = pnu.NextSet(n + 1) {
				// remaining code like BK1
				R.Copy(r2)
				r2.Set(n)
				p2.ClearAll()
				x2.ClearAll()
				for _, to := range g.AdjacencyList[n] {
					if P.Test(uint(to)) {
						p2.Set(uint(to))
					}
					if X.Test(uint(to)) {
						x2.Set(uint(to))
					}
				}
				if !f(r2, p2, x2) {
					return false
				}
				P.SetTo(n, false)
				X.Set(n)
			}
		case X.None():
			var n uint
			n--
			c := make([]NI, R.Count())
			for i := range c {
				n, _ = R.NextSet(n + 1)
				c[i] = NI(n)
			}
			if !emit(c) {
				return false
			}
		}
		return true
	}
	R := bitset.New(uint(len(g.AdjacencyList)))
	P := bitset.New(uint(len(g.AdjacencyList))).Complement()
	X := bitset.New(uint(len(g.AdjacencyList)))
	f(R, P, X)
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
// The method calls the emit argument for each maximal clique in g, as long
// as emit returns true.  If emit returns false, BronKerbosch1 returns
// immediately.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also simpler variants BronKerbosch1 and BronKerbosch2.
func (g Undirected) BronKerbosch3(pivot func(P, X *bitset.BitSet) int, emit func([]NI) bool) {
	var f func(R, P, X *bitset.BitSet) bool
	f = func(R, P, X *bitset.BitSet) bool {
		switch {
		case P.Any():
			r2 := bitset.New(uint(len(g.AdjacencyList)))
			p2 := bitset.New(uint(len(g.AdjacencyList)))
			x2 := bitset.New(uint(len(g.AdjacencyList)))
			// compute P \ N(u).  next 5 lines are only difference from BK1
			pnu := P.Clone()
			for _, to := range g.AdjacencyList[pivot(P, X)] {
				pnu.SetTo(uint(to), false)
			}
			for n, ok := pnu.NextSet(0); ok; n, ok = pnu.NextSet(n + 1) {
				// remaining code like BK1
				R.Copy(r2)
				r2.Set(n)
				p2.ClearAll()
				x2.ClearAll()
				for _, to := range g.AdjacencyList[n] {
					if P.Test(uint(to)) {
						p2.Set(uint(to))
					}
					if X.Test(uint(to)) {
						x2.Set(uint(to))
					}
				}
				if !f(r2, p2, x2) {
					return false
				}
				P.SetTo(n, false)
				X.Set(n)
			}
		case X.None():
			var n uint
			n--
			c := make([]NI, R.Count())
			for i := range c {
				n, _ = R.NextSet(n + 1)
				c[i] = NI(n)
			}
			if !emit(c) {
				return false
			}
		}
		return true
	}
	R := bitset.New(uint(len(g.AdjacencyList)))
	P := bitset.New(uint(len(g.AdjacencyList))).Complement()
	X := bitset.New(uint(len(g.AdjacencyList)))
	// code above same as BK2
	// code below new to BK3
	_, ord, _ := g.Degeneracy()
	p2 := bitset.New(uint(len(g.AdjacencyList)))
	x2 := bitset.New(uint(len(g.AdjacencyList)))
	for _, n := range ord {
		R.Set(uint(n))
		p2.ClearAll()
		x2.ClearAll()
		for _, to := range g.AdjacencyList[n] {
			if P.Test(uint(to)) {
				p2.Set(uint(to))
			}
			if X.Test(uint(to)) {
				x2.Set(uint(to))
			}
		}
		if !f(R, p2, x2) {
			return
		}
		R.SetTo(uint(n), false)
		P.SetTo(uint(n), false)
		X.Set(uint(n))
	}
}

// ConnectedComponentBits returns a function that iterates over connected
// components of g, returning a member bitmap for each.
//
// Each call of the returned function returns the order (number of nodes)
// and bits of a connected component.  The returned function returns zeros
// after returning all connected components.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps, which has lighter weight return values.
func (g Undirected) ConnectedComponentBits() func() (order int, bits big.Int) {
	var vg big.Int  // nodes visited in graph
	var vc *big.Int // nodes visited in current component
	var nc int
	var df func(NI)
	df = func(n NI) {
		vg.SetBit(&vg, int(n), 1)
		vc.SetBit(vc, int(n), 1)
		nc++
		for _, nb := range g.AdjacencyList[n] {
			if vg.Bit(int(nb)) == 0 {
				df(nb)
			}
		}
		return
	}
	var n NI
	return func() (o int, bits big.Int) {
		for ; n < NI(len(g.AdjacencyList)); n++ {
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

// ConnectedComponentLists returns a function that iterates over connected
// components of g, returning the member list of each.
//
// Each call of the returned function returns a node list of a connected
// component.  The returned function returns nil after returning all connected
// components.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps, which has lighter weight return values.
func (g Undirected) ConnectedComponentLists() func() []NI {
	var vg big.Int // nodes visited in graph
	var m []NI     // members of current component
	var df func(NI)
	df = func(n NI) {
		vg.SetBit(&vg, int(n), 1)
		m = append(m, n)
		for _, nb := range g.AdjacencyList[n] {
			if vg.Bit(int(nb)) == 0 {
				df(nb)
			}
		}
		return
	}
	var n NI
	return func() []NI {
		for ; n < NI(len(g.AdjacencyList)); n++ {
			if vg.Bit(int(n)) == 0 {
				m = nil
				df(n)
				return m
			}
		}
		return nil
	}
}

// ConnectedComponentReps returns a representative node from each connected
// component of g.
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
// collect component members in a single traversal, and IsConnected which
// is an even simpler boolean test.
func (g Undirected) ConnectedComponentReps() (reps []NI, orders []int) {
	var c big.Int
	var o int
	var df func(NI)
	df = func(n NI) {
		c.SetBit(&c, int(n), 1)
		o++
		for _, nb := range g.AdjacencyList[n] {
			if c.Bit(int(nb)) == 0 {
				df(nb)
			}
		}
		return
	}
	for n := range g.AdjacencyList {
		if c.Bit(n) == 0 {
			reps = append(reps, NI(n))
			o = 0
			df(NI(n))
			orders = append(orders, o)
		}
	}
	return
}

// Copy makes a deep copy of g.
// Copy also computes the arc size ma, the number of arcs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g Undirected) Copy() (c Undirected, ma int) {
	l, s := g.AdjacencyList.Copy()
	return Undirected{l}, s
}

// Degeneracy computes k-degeneracy, vertex ordering and k-cores.
//
// See Wikipedia https://en.wikipedia.org/wiki/Degeneracy_(graph_theory)
//
// There are equivalent labeled and unlabeled versions of this method.
func (g Undirected) Degeneracy() (k int, ordering []NI, cores []int) {
	// WP algorithm
	ordering = make([]NI, len(g.AdjacencyList))
	var L big.Int
	d := make([]int, len(g.AdjacencyList))
	var D [][]NI
	for v, nb := range g.AdjacencyList {
		dv := len(nb)
		d[v] = dv
		for len(D) <= dv {
			D = append(D, nil)
		}
		D[dv] = append(D[dv], NI(v))
	}
	for ox := range g.AdjacencyList {
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
		for _, nb := range g.AdjacencyList[v] {
			if L.Bit(int(nb)) == 1 {
				continue
			}
			dn := d[nb]  // old number of neighbors of nb
			Ddn := D[dn] // nb is in this list
			// remove it from the list
			for wx, w := range Ddn {
				if w == nb {
					last := len(Ddn) - 1
					Ddn[wx], Ddn[last] = Ddn[last], Ddn[wx]
					D[dn] = Ddn[:last]
				}
			}
			dn-- // new number of neighbors
			d[nb] = dn
			// re--add it to it's new list
			D[dn] = append(D[dn], nb)
		}
	}
	cores[k] = len(ordering)
	return
}

// Degree for undirected graphs, returns the degree of a node.
//
// The degree of a node in an undirected graph is the number of incident
// edges, where loops count twice.
//
// If g is known to be loop-free, the result is simply equivalent to len(g[n]).
// See handshaking lemma example at AdjacencyList.ArcSize.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g Undirected) Degree(n NI) int {
	to := g.AdjacencyList[n]
	d := len(to) // just "out" degree,
	for _, to := range to {
		if to == n {
			d++ // except loops count twice
		}
	}
	return d
}

// IsConnected tests if an undirected graph is a single connected component.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps for a method returning more information.
func (g Undirected) IsConnected() bool {
	if len(g.AdjacencyList) == 0 {
		return true
	}
	var b big.Int
	OneBits(&b, len(g.AdjacencyList))
	var df func(int)
	df = func(n int) {
		b.SetBit(&b, n, 0)
		for _, to := range g.AdjacencyList[n] {
			to := int(to)
			if b.Bit(to) == 1 {
				df(to)
			}
		}
	}
	df(0)
	return b.BitLen() == 0
}

// IsTree identifies trees in undirected graphs.
//
// IsTree returns true if the connected component
// containing argument root is a tree.  It does not validate
// that the entire graph is a tree.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g Undirected) IsTree(root NI) bool {
	var v big.Int
	var df func(NI, NI) bool
	df = func(fr, n NI) bool {
		if v.Bit(int(n)) == 1 {
			return false
		}
		v.SetBit(&v, int(n), 1)
		for _, to := range g.AdjacencyList[n] {
			if to != fr && !df(n, to) {
				return false
			}
		}
		return true
	}
	v.SetBit(&v, int(root), 1)
	for _, to := range g.AdjacencyList[root] {
		if !df(root, to) {
			return false
		}
	}
	return true
}
