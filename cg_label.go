// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import "math/big"

// cg_adj.go is code generated from cg_label.go by directive in graph.go.
// Editing cg_label.go is okay.
// DO NOT EDIT cg_adj.go.

// ArcSize returns the number of arcs in g.
//
// Note that for an undirected graph witout loops, the number of edges --
// the traditional meaning of graph size -- will be m/2.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) ArcSize() (m int) {
	for _, to := range g {
		m += len(to)
	}
	return
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
func (g LabeledAdjacencyList) Bipartite(n int) (b bool, c1, c2 *big.Int, oc []int) {
	c1 = &big.Int{}
	c2 = &big.Int{}
	b = true
	var open bool
	var df func(n int, c1, c2 *big.Int)
	df = func(n int, c1, c2 *big.Int) {
		c1.SetBit(c1, n, 1)
		for _, nb := range g[n] {
			if c1.Bit(nb.To) == 1 {
				b = false
				oc = []int{nb.To, n}
				open = true
				return
			}
			if c2.Bit(nb.To) == 1 {
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
func (g LabeledAdjacencyList) BoundsOk() (ok bool, fr int, to Half) {
	for fr, to := range g {
		for _, to := range to {
			if to.To < 0 || to.To >= len(g) {
				return false, fr, to
			}
		}
	}
	return true, -1, to
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
	var df func(int)
	df = func(n int) {
		vg.SetBit(&vg, n, 1)
		vc.SetBit(vc, n, 1)
		nc++
		for _, nb := range g[n] {
			if vg.Bit(nb.To) == 0 {
				df(nb.To)
			}
		}
		return
	}
	n := 0
	return func() (o int, bits big.Int) {
		for ; n < len(g); n++ {
			if vg.Bit(n) == 0 {
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
func (g LabeledAdjacencyList) ConnectedComponentLists() func() []int {
	var vg big.Int // nodes visited in graph
	var m []int    // members of current component
	var df func(int)
	df = func(n int) {
		vg.SetBit(&vg, n, 1)
		m = append(m, n)
		for _, nb := range g[n] {
			if vg.Bit(nb.To) == 0 {
				df(nb.To)
			}
		}
		return
	}
	n := 0
	return func() []int {
		for ; n < len(g); n++ {
			if vg.Bit(n) == 0 {
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
func (g LabeledAdjacencyList) ConnectedComponentReps() (reps, orders []int) {
	var c big.Int
	var o int
	var df func(int)
	df = func(n int) {
		c.SetBit(&c, n, 1)
		o++
		for _, nb := range g[n] {
			if c.Bit(nb.To) == 0 {
				df(nb.To)
			}
		}
		return
	}
	for n := range g {
		if c.Bit(n) == 0 {
			reps = append(reps, n)
			o = 0
			df(n)
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
func (g LabeledAdjacencyList) Cyclic() bool {
	var c bool
	var temp, perm big.Int
	var df func(int)
	df = func(n int) {
		switch {
		case temp.Bit(n) == 1:
			c = true
			return
		case perm.Bit(n) == 1:
			return
		}
		temp.SetBit(&temp, n, 1)
		for _, nb := range g[n] {
			df(nb.To)
			if c {
				return
			}
		}
		temp.SetBit(&temp, n, 0)
		perm.SetBit(&perm, n, 1)
	}
	for n := range g {
		if perm.Bit(n) == 1 {
			continue
		}
		if df(n); c {
			break
		}
	}
	return c
}

//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) Degeneracy() (k int, ord []int, cores []int) {
	ord = make([]int, len(g))
	var L big.Int
	d := make([]int, len(g))
	var D [][]int
	for v, nb := range g {
		dv := len(nb)
		d[v] = dv
		for len(D) <= dv {
			D = append(D, nil)
		}
		D[dv] = append(D[dv], v)
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
		ord[ox] = v
		L.SetBit(&L, v, 1)
		D[i] = Di[:last]
		// move neighbors
		for _, nb := range g[v] {
			if L.Bit(nb.To) == 1 {
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
	cores[k] = len(ord)
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
func (g LabeledAdjacencyList) DepthFirst(start int, bm *big.Int, v Visitor) (ok bool) {
	if bm == nil {
		if v == nil {
			return false
		}
		bm = new(big.Int)
	}
	ok = true
	var df func(n int)
	df = func(n int) {
		if bm.Bit(n) == 1 {
			return
		}
		bm.SetBit(bm, n, 1)
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
	leaves.Sub(leaves.Lsh(one, uint(len(g))), one)
	// iterate over arcs, setting from pointers and and marking non-leaves.
	for fr, to := range g {
		for _, to := range to {
			paths[to.To].From = fr
			leaves.SetBit(&leaves, fr, 0)
		}
	}
	// f to set path lengths
	var leng func(int) int
	leng = func(n int) int {
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
			if l := leng(i); l > maxLen {
				maxLen = l
			}
		}
	}
	return FromList{paths, leaves, maxLen}
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

// IsTreeDirected identifies trees in directed graphs.
//
// IsTreeDirected returns true if the subgraph reachable from
// root is a tree.  It does not validate that the entire graph is a tree.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) IsTreeDirected(root int) bool {
	var v big.Int
	var df func(int) bool
	df = func(n int) bool {
		if v.Bit(n) == 1 {
			return false
		}
		v.SetBit(&v, n, 1)
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
func (g LabeledAdjacencyList) IsTreeUndirected(root int) bool {
	var v big.Int
	var df func(int, int) bool
	df = func(fr, n int) bool {
		if v.Bit(n) == 1 {
			return false
		}
		v.SetBit(&v, n, 1)
		for _, to := range g[n] {
			if to.To != fr && !df(n, to.To) {
				return false
			}
		}
		return true
	}
	v.SetBit(&v, root, 1)
	for _, to := range g[root] {
		if !df(root, to.To) {
			return false
		}
	}
	return true
}

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
func (g LabeledAdjacencyList) Tarjan() (scc [][]int) {
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
	var S []int
	var sc func(int)
	sc = func(n int) {
		index[n] = x
		indexed.SetBit(&indexed, n, 1)
		lowlink[n] = x
		x++
		S = append(S, n)
		stacked.SetBit(&stacked, n, 1)
		for _, nb := range g[n] {
			if indexed.Bit(nb.To) == 0 {
				sc(nb.To)
				if lowlink[nb.To] < lowlink[n] {
					lowlink[n] = lowlink[nb.To]
				}
			} else if stacked.Bit(nb.To) == 1 {
				if index[nb.To] < lowlink[n] {
					lowlink[n] = index[nb.To]
				}
			}
		}
		if lowlink[n] == index[n] {
			var c []int
			for {
				last := len(S) - 1
				w := S[last]
				S = S[:last]
				stacked.SetBit(&stacked, w, 0)
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
			sc(n)
		}
	}
	return
}

// TarjanForward returns strongly connected components.
//
// It returns components in the reverse order of Tarjan, for situations
// where a forward topological ordering is easier.
func (g LabeledAdjacencyList) TarjanForward() (scc [][]int) {
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
func (g LabeledAdjacencyList) TarjanCondensation() (scc [][]int, cd AdjacencyList) {
	scc = g.TarjanForward()
	cd = make(AdjacencyList, len(scc)) // return value
	cond := make([]int, len(g))        // mapping from g node to cd node
	for cn := len(scc) - 1; cn >= 0; cn-- {
		c := scc[cn]
		for _, n := range c {
			cond[n] = cn // map g node to cd node
		}
		var tos []int // list of 'to' nodes
		var m big.Int // tos map
		m.SetBit(&m, cn, 1)
		for _, n := range c {
			for _, to := range g[n] {
				if ct := cond[to.To]; m.Bit(ct) == 0 {
					m.SetBit(&m, ct, 1)
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
func (g LabeledAdjacencyList) Topological() (ordering, cycle []int) {
	ordering = make([]int, len(g))
	i := len(ordering)
	var temp, perm big.Int
	var cycleFound bool
	var cycleStart int
	var df func(int)
	df = func(n int) {
		switch {
		case temp.Bit(n) == 1:
			cycleFound = true
			cycleStart = n
			return
		case perm.Bit(n) == 1:
			return
		}
		temp.SetBit(&temp, n, 1)
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
		temp.SetBit(&temp, n, 0)
		perm.SetBit(&perm, n, 1)
		i--
		ordering[i] = n
	}
	for n := range g {
		if perm.Bit(n) == 1 {
			continue
		}
		df(n)
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
func (g LabeledAdjacencyList) TopologicalKahn(tr AdjacencyList) (ordering, cycle []int) {
	// code follows Wikipedia pseudocode.
	var L, S []int
	// rem for "remaining edges," this function makes a local copy of the
	// in-degrees and consumes that instead of consuming an input.
	rem := make([]int, len(g))
	for n, fr := range tr {
		if len(fr) == 0 {
			// accumulate "set of all nodes with no incoming edges"
			S = append(S, n)
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
					cycle = append(cycle, c)
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
