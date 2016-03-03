// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"math/big"
	"math/rand"
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

// BoundsOk validates that all arcs in g stay within the slice bounds of g.
//
// BoundsOk returns true when no arcs point outside the bounds of g.
// Otherwise it returns false and an example arc that points outside of g.
//
// Most methods of this package assume the BoundsOk condition and may
// panic when they encounter an arc pointing outside of the graph.  This
// function can be used to validate a graph when the BoundsOk condition
// is unknown.
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

// Copy makes a deep copy of g.
// Copy also computes the arc size ma, the number of arcs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) Copy() (c LabeledAdjacencyList, ma int) {
	c = make(LabeledAdjacencyList, len(g))
	for n, to := range g {
		c[n] = append([]Half{}, to...)
		ma += len(to)
	}
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

// DepthFirstRandom traverses a graph depth first, but following arcs in
// random order among arcs from a single node.
//
// Usage is otherwise like the DepthFirst method.  See DepthFirst.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) DepthFirstRandom(start NI, bm *big.Int, v Visitor, r *rand.Rand) (ok bool) {
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
		to := g[n]
		for _, i := range r.Perm(len(to)) {
			df(to[i].To)
		}
	}
	df(start)
	return
}

// HasArc returns true if g has any arc from node fr to node to.
//
// Also returned is the index within the slice of arcs from node fr.
// If no arc from fr to to is present, HasArc returns false, -1.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) HasArc(fr, to NI) (bool, int) {
	for x, h := range g[fr] {
		if h.To == to {
			return true, x
		}
	}
	return false, -1
}

// HasLoop identifies if a graph contains a loop, an arc that leads from a
// a node back to the same node.
//
// If the graph has a loop, the result is an example node that has a loop.
//
// If g contains a loop, the method returns true and an example of a node
// with a loop.  If there are no loops in g, the method returns false, -1.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) HasLoop() (bool, NI) {
	for fr, to := range g {
		for _, to := range to {
			if NI(fr) == to.To {
				return true, to.To
			}
		}
	}
	return false, -1
}

// HasParallelMap identifies if a graph contains parallel arcs, multiple arcs
// that lead from a node to the same node.
//
// If the graph has parallel arcs, the method returns true and
// results fr and to represent an example where there are parallel arcs
// from node fr to node to.
//
// If there are no parallel arcs, the method returns false, -1 -1.
//
// Multiple loops on a node count as parallel arcs.
//
// "Map" in the method name indicates that a Go map is used to detect parallel
// arcs.  Compared to method HasParallelSort, this gives better asymtotic
// performance for large dense graphs but may have increased overhead for
// small or sparse graphs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledAdjacencyList) HasParallelMap() (has bool, fr, to NI) {
	for n, to := range g {
		if len(to) == 0 {
			continue
		}
		m := map[NI]struct{}{}
		for _, to := range to {
			if _, ok := m[to.To]; ok {
				return true, NI(n), to.To
			}
			m[to.To] = struct{}{}
		}
	}
	return false, -1, -1
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
	if lp, n := g.HasLoop(); lp {
		return false, n
	}
	if pa, n, _ := g.HasParallelSort(); pa {
		return false, n
	}
	return true, -1
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
