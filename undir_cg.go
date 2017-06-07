// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"github.com/soniakeys/bits"
)

// undir_RO.go is code generated from undir_cg.go by directives in graph.go.
// Editing undir_cg.go is okay.  It is the code generation source.
// DO NOT EDIT undir_RO.go.
// The RO means read only and it is upper case RO to slow you down a bit
// in case you start to edit the file.

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
func (g LabeledUndirected) Bipartite(n NI) (b bool, c1, c2 bits.Bits, oc []NI) {
	a := g.LabeledAdjacencyList
	c1 = bits.New(len(a))
	c2 = bits.New(len(a))
	b = true
	var open bool
	var df func(n NI, c1, c2 *bits.Bits)
	df = func(n NI, c1, c2 *bits.Bits) {
		c1.SetBit(int(n), 1)
		for _, nb := range a[n] {
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
	df(n, &c1, &c2)
	if b {
		return b, c1, c2, nil
	}
	return b, bits.Bits{}, bits.Bits{}, oc
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
func (g LabeledUndirected) BronKerbosch1(emit func(bits.Bits) bool) {
	a := g.LabeledAdjacencyList
	var f func(R, P, X bits.Bits) bool
	f = func(R, P, X bits.Bits) bool {
		switch {
		case !P.Zero():
			r2 := bits.New(len(a))
			p2 := bits.New(len(a))
			x2 := bits.New(len(a))
			pf := func(n int) bool {
				r2.Set(R)
				r2.SetBit(n, 1)
				p2.ClearAll()
				x2.ClearAll()
				for _, to := range a[n] {
					if P.Bit(int(to.To)) == 1 {
						p2.SetBit(int(to.To), 1)
					}
					if X.Bit(int(to.To)) == 1 {
						x2.SetBit(int(to.To), 1)
					}
				}
				if !f(r2, p2, x2) {
					return false
				}
				P.SetBit(n, 0)
				X.SetBit(n, 1)
				return true
			}
			if !P.IterateOnes(pf) {
				return false
			}
		case X.Zero():
			return emit(R)
		}
		return true
	}
	var R, P, X bits.Bits
	R = bits.New(len(a))
	P = bits.New(len(a))
	X = bits.New(len(a))
	P.SetAll()
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
func (g LabeledUndirected) BKPivotMaxDegree(P, X bits.Bits) (p NI) {
	// choose pivot u as highest degree node from P or X
	a := g.LabeledAdjacencyList
	maxDeg := -1
	P.IterateOnes(func(n int) bool { // scan P
		if d := len(a[n]); d > maxDeg {
			p = NI(n)
			maxDeg = d
		}
		return true
	})
	X.IterateOnes(func(n int) bool { // scan X
		if d := len(a[n]); d > maxDeg {
			p = NI(n)
			maxDeg = d
		}
		return true
	})
	return
}

// BKPivotMinP is a strategy for BronKerbosch methods.
//
// To use it, take the method value (see golang.org/ref/spec#Method_values)
// and pass it as the argument to BronKerbosch2 or 3.
//
// The strategy is to simply pick the first node in P.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) BKPivotMinP(P, X bits.Bits) NI {
	return NI(P.OneFrom(0))
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
func (g LabeledUndirected) BronKerbosch2(pivot func(P, X bits.Bits) NI, emit func(bits.Bits) bool) {
	a := g.LabeledAdjacencyList
	var f func(R, P, X bits.Bits) bool
	f = func(R, P, X bits.Bits) bool {
		switch {
		case !P.Zero():
			r2 := bits.New(len(a))
			p2 := bits.New(len(a))
			x2 := bits.New(len(a))
			pnu := bits.New(len(a))
			// compute P \ N(u).  next 5 lines are only difference from BK1
			pnu.Set(P)
			for _, to := range a[pivot(P, X)] {
				pnu.SetBit(int(to.To), 0)
			}
			// remaining code like BK1
			pf := func(n int) bool {
				r2.Set(R)
				r2.SetBit(n, 1)
				p2.ClearAll()
				x2.ClearAll()
				for _, to := range a[n] {
					if P.Bit(int(to.To)) == 1 {
						p2.SetBit(int(to.To), 1)
					}
					if X.Bit(int(to.To)) == 1 {
						x2.SetBit(int(to.To), 1)
					}
				}
				if !f(r2, p2, x2) {
					return false
				}
				P.SetBit(n, 0)
				X.SetBit(n, 1)
				return true
			}
			if !pnu.IterateOnes(pf) {
				return false
			}
		case X.Zero():
			return emit(R)
		}
		return true
	}
	R := bits.New(len(a))
	P := bits.New(len(a))
	X := bits.New(len(a))
	P.SetAll()
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
func (g LabeledUndirected) BronKerbosch3(pivot func(P, X bits.Bits) NI, emit func(bits.Bits) bool) {
	a := g.LabeledAdjacencyList
	var f func(R, P, X bits.Bits) bool
	f = func(R, P, X bits.Bits) bool {
		switch {
		case !P.Zero():
			r2 := bits.New(len(a))
			p2 := bits.New(len(a))
			x2 := bits.New(len(a))
			pnu := bits.New(len(a))
			// compute P \ N(u).  next lines are only difference from BK1
			pnu.Set(P)
			for _, to := range a[pivot(P, X)] {
				pnu.SetBit(int(to.To), 0)
			}
			// remaining code like BK2
			pf := func(n int) bool {
				r2.Set(R)
				r2.SetBit(n, 1)
				p2.ClearAll()
				x2.ClearAll()
				for _, to := range a[n] {
					if P.Bit(int(to.To)) == 1 {
						p2.SetBit(int(to.To), 1)
					}
					if X.Bit(int(to.To)) == 1 {
						x2.SetBit(int(to.To), 1)
					}
				}
				if !f(r2, p2, x2) {
					return false
				}
				P.SetBit(n, 0)
				X.SetBit(n, 1)
				return true
			}
			if !pnu.IterateOnes(pf) {
				return false
			}
		case X.Zero():
			return emit(R)
		}
		return true
	}
	R := bits.New(len(a))
	P := bits.New(len(a))
	X := bits.New(len(a))
	P.SetAll()
	// code above same as BK2
	// code below new to BK3
	_, ord, _ := g.Degeneracy()
	p2 := bits.New(len(a))
	x2 := bits.New(len(a))
	for _, n := range ord {
		R.SetBit(int(n), 1)
		p2.ClearAll()
		x2.ClearAll()
		for _, to := range a[n] {
			if P.Bit(int(to.To)) == 1 {
				p2.SetBit(int(to.To), 1)
			}
			if X.Bit(int(to.To)) == 1 {
				x2.SetBit(int(to.To), 1)
			}
		}
		if !f(R, p2, x2) {
			return
		}
		R.SetBit(int(n), 0)
		P.SetBit(int(n), 0)
		X.SetBit(int(n), 1)
	}
}

// ConnectedComponentbits.Bits returns a function that iterates over connected
// components of g, returning a member bitmap for each.
//
// Each call of the returned function returns the order, arc size,
// and bits of a connected component.  The underlying bits allocation is
// the same for each call and is overwritten on subsequent calls.  Use or
// save the bits before calling the function again.  The function returns
// zeros after returning all connected components.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps, which has lighter weight return values.
func (g LabeledUndirected) ConnectedComponentBits() func() (order, arcSize int, bits bits.Bits) {
	a := g.LabeledAdjacencyList
	vg := bits.New(len(a)) // nodes visited in graph
	vc := bits.New(len(a)) // nodes visited in current component
	var order, arcSize int
	var df func(NI)
	df = func(n NI) {
		vg.SetBit(int(n), 1)
		vc.SetBit(int(n), 1)
		order++
		arcSize += len(a[n])
		for _, nb := range a[n] {
			if vg.Bit(int(nb.To)) == 0 {
				df(nb.To)
			}
		}
		return
	}
	var n int
	return func() (o, ma int, b bits.Bits) {
		for ; n < len(a); n++ {
			if vg.Bit(n) == 0 {
				vc.ClearAll()
				order, arcSize = 0, 0
				df(NI(n))
				return order, arcSize, vc
			}
		}
		return // return zeros signalling no more components
	}
}

// ConnectedComponentLists returns a function that iterates over connected
// components of g, returning the member list of each.
//
// Each call of the returned function returns a node list of a connected
// component and the arc size of the component.  The returned function returns
// nil, 0 after returning all connected components.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps, which has lighter weight return values.
func (g LabeledUndirected) ConnectedComponentLists() func() (nodes []NI, arcSize int) {
	a := g.LabeledAdjacencyList
	vg := bits.New(len(a)) // nodes visited in graph
	var l []NI             // accumulated node list of current component
	var ma int             // accumulated arc size of current component
	var df func(NI)
	df = func(n NI) {
		vg.SetBit(int(n), 1)
		l = append(l, n)
		ma += len(a[n])
		for _, nb := range a[n] {
			if vg.Bit(int(nb.To)) == 0 {
				df(nb.To)
			}
		}
		return
	}
	var n int
	return func() ([]NI, int) {
		for ; n < len(a); n++ {
			if vg.Bit(n) == 0 {
				l, ma = nil, 0
				df(NI(n))
				return l, ma
			}
		}
		return nil, 0
	}
}

// ConnectedComponentReps returns a representative node from each connected
// component of g.
//
// Returned is a slice with a single representative node from each connected
// component and also parallel slices with the orders and arc sizes
// in the corresponding components.
//
// This is fairly minimal information describing connected components.
// From a representative node, other nodes in the component can be reached
// by depth first traversal for example.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentbits.Bits and ConnectedComponentLists which can
// collect component members in a single traversal, and IsConnected which
// is an even simpler boolean test.
func (g LabeledUndirected) ConnectedComponentReps() (reps []NI, orders, arcSizes []int) {
	a := g.LabeledAdjacencyList
	c := bits.New(len(a))
	var o, ma int
	var df func(NI)
	df = func(n NI) {
		c.SetBit(int(n), 1)
		o++
		ma += len(a[n])
		for _, nb := range a[n] {
			if c.Bit(int(nb.To)) == 0 {
				df(nb.To)
			}
		}
		return
	}
	for n := range a {
		if c.Bit(n) == 0 {
			o, ma = 0, 0
			df(NI(n))
			reps = append(reps, NI(n))
			orders = append(orders, o)
			arcSizes = append(arcSizes, ma)
		}
	}
	return
}

// Copy makes a deep copy of g.
// Copy also computes the arc size ma, the number of arcs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) Copy() (c LabeledUndirected, ma int) {
	l, s := g.LabeledAdjacencyList.Copy()
	return LabeledUndirected{l}, s
}

// Degeneracy computes k-degeneracy, vertex ordering and k-cores.
//
// See Wikipedia https://en.wikipedia.org/wiki/Degeneracy_(graph_theory)
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) Degeneracy() (k int, ordering []NI, cores []int) {
	a := g.LabeledAdjacencyList
	// WP algorithm
	ordering = make([]NI, len(a))
	L := bits.New(len(a))
	d := make([]int, len(a))
	var D [][]NI
	for v, nb := range a {
		dv := len(nb)
		d[v] = dv
		for len(D) <= dv {
			D = append(D, nil)
		}
		D[dv] = append(D[dv], NI(v))
	}
	for ox := range a {
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
		L.SetBit(int(v), 1)
		D[i] = Di[:last]
		// move neighbors
		for _, nb := range a[v] {
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

// Degree for undirected graphs, returns the degree of a node.
//
// The degree of a node in an undirected graph is the number of incident
// edges, where loops count twice.
//
// If g is known to be loop-free, the result is simply equivalent to len(g[n]).
// See handshaking lemma example at AdjacencyList.ArcSize.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) Degree(n NI) int {
	to := g.LabeledAdjacencyList[n]
	d := len(to) // just "out" degree,
	for _, to := range to {
		if to.To == n {
			d++ // except loops count twice
		}
	}
	return d
}

// Density returns density for a simple graph.
//
// See also Density function.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) Density() float64 {
	return Density(g.Order(), g.Size())
}

// FromList constructs a FromList representing the tree reachable from
// the given root.
//
// The connected component containing root should represent a simple graph,
// connected as a tree.
//
// For nodes connected as a tree, the Path member of the returned FromList
// will be populated with both From and Len values.  The MaxLen member will be
// set but Leaves will be left a zero value.  Return value cycle will be -1.
//
// If the connected component containing root is not connected as a tree,
// a cycle will be detected.  The returned FromList will be a zero value and
// return value cycle will be a node involved in the cycle.
//
// Loops and parallel edges will be detected as cycles, however only in the
// connected component containing root.  If g is not fully connected, nodes
// not reachable from root will have PathEnd values of {From: -1, Len: 0}.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) FromList(root NI) (f FromList, cycle NI) {
	p := make([]PathEnd, g.Order())
	for i := range p {
		p[i].From = -1
	}
	ml := 0
	var df func(NI, NI) bool
	df = func(fr, n NI) bool {
		l := p[n].Len + 1
		for _, to := range g.LabeledAdjacencyList[n] {
			if to.To == fr {
				continue
			}
			if p[to.To].Len > 0 {
				cycle = to.To
				return false
			}
			p[to.To] = PathEnd{From: n, Len: l}
			if l > ml {
				ml = l
			}
			if !df(n, to.To) {
				return false
			}
		}
		return true
	}
	p[root].Len = 1
	if !df(-1, root) {
		return
	}
	return FromList{Paths: p, MaxLen: ml}, -1
}

// IsConnected tests if an undirected graph is a single connected component.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps for a method returning more information.
func (g LabeledUndirected) IsConnected() bool {
	a := g.LabeledAdjacencyList
	if len(a) == 0 {
		return true
	}
	b := bits.New(len(a))
	b.SetAll()
	var df func(NI)
	df = func(n NI) {
		b.SetBit(int(n), 0)
		for _, to := range a[n] {
			if b.Bit(int(to.To)) == 1 {
				df(to.To)
			}
		}
	}
	df(0)
	return b.Zero()
}

// IsTree identifies trees in undirected graphs.
//
// Return value isTree is true if the connected component reachable from root
// is a tree.  Further, return value allTree is true if the entire graph g is
// connected.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) IsTree(root NI) (isTree, allTree bool) {
	a := g.LabeledAdjacencyList
	v := bits.New(len(a))
	v.SetAll()
	var df func(NI, NI) bool
	df = func(fr, n NI) bool {
		if v.Bit(int(n)) == 0 {
			return false
		}
		v.SetBit(int(n), 0)
		for _, to := range a[n] {
			if to.To != fr && !df(n, to.To) {
				return false
			}
		}
		return true
	}
	v.SetBit(int(root), 0)
	for _, to := range a[root] {
		if !df(root, to.To) {
			return false, false
		}
	}
	return true, v.Zero()
}

// Size returns the number of edges in g.
//
// See also ArcSize and AnyLoop.
func (g LabeledUndirected) Size() int {
	m2 := 0
	for fr, to := range g.LabeledAdjacencyList {
		m2 += len(to)
		for _, to := range to {
			if to.To == NI(fr) {
				m2++
			}
		}
	}
	return m2 / 2
}
