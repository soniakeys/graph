// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// dir.go
//
// Methods specific to directed graphs.
// Doc for each method should specifically say directed.

package graph

import (
	"errors"
	"math/big"
)

// Transpose, for directed graphs, constructs a new adjacency list that is
// the transpose of g.
//
// For every arc from->to of g, the result will have an arc to->from.
// Transpose also counts arcs as it traverses and returns m the number of arcs
// in g (equal to the number of arcs in the result.)
func (g AdjacencyList) Transpose() (t AdjacencyList, m int) {
	t = make([][]int, len(g))
	for n, nbs := range g {
		for _, nb := range nbs {
			t[nb] = append(t[nb], n)
			m++
		}
	}
	return
}

// Cyclic, for directed graphs, determines if g contains cycles.
//
// Cyclic returns true if g contains at least one cycle.
// Cyclic returns false if g is acyclic.
func (g AdjacencyList) Cyclic() bool {
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
			df(nb)
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

// Topological, for directed acyclic graphs, computes a topological sort of g.
//
// For an acyclic graph, return value order is a permutation of node numbers
// in topologically sorted order and cycle will be nil.  If the graph is found
// to be cyclic, order will be nil and cycle will be the path of a found cycle.
func (g AdjacencyList) Topological() (order, cycle []int) {
	order = make([]int, len(g))
	i := len(order)
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
			df(nb)
			if cycleFound {
				if cycleStart >= 0 {
					cycle = append(cycle, n)
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
		order[i] = n
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
	return order, nil
}

// Tarjan identifies strongly connected components in a directed graph using
// Tarjan's algorithm.
//
// Returned is a list of components, each component is a list of nodes.
func (g AdjacencyList) Tarjan() (scc [][]int) {
	// straight from WP
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
			if indexed.Bit(nb) == 0 {
				sc(nb)
				if lowlink[nb] < lowlink[n] {
					lowlink[n] = lowlink[nb]
				}
			} else if stacked.Bit(nb) == 1 {
				if index[nb] < lowlink[n] {
					lowlink[n] = index[nb]
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
	return scc
}

// InDegree computes the in-degree of each node in g
func (g AdjacencyList) InDegree() []int {
	ind := make([]int, len(g))
	for _, nbs := range g {
		for _, nb := range nbs {
			ind[nb]++
		}
	}
	return ind
}

// Balanced returns true if for every node in g, in-degree equals out-degree.
func (g AdjacencyList) Balanced() bool {
	for n, in := range g.InDegree() {
		if in != len(g[n]) {
			return false
		}
	}
	return true
}

// EulerianCycle finds an Eulerian cycle in a directed multigraph.
//
// * If g has no nodes, result is nil, nil.
//
// * If g is Eulerian, result is an Eulerian cycle with err = nil.
// The cycle result is a list of nodes, where the first and last
// nodes are the same.
//
// * Otherwise, result is nil, error
//
// Internally, EulerianCycle copies the entire graph g.
// See EulerianCycleD for a more space efficient version.
func (g AdjacencyList) EulerianCycle() ([]int, error) {
	c, m := g.Copy()
	return c.EulerianCycleD(m)
}

// EulerianCycleD finds an Eulerian cycle in a directed multigraph.
//
// EulerianCycleD is destructive on its receiver g.  See EulerianCycle for
// a non-destructive version.
//
// Argument m must be the correct size, or number of arcs in g.
//
// * If g has no nodes, result is nil, nil.
//
// * If g is Eulerian, result is an Eulerian cycle with err = nil.
// The cycle result is a list of nodes, where the first and last
// nodes are the same.
//
// * Otherwise, result is nil, error
func (g AdjacencyList) EulerianCycleD(m int) ([]int, error) {
	if len(g) == 0 {
		return nil, nil
	}
	p := make([]int, m+1)
	for s := []int{0}; len(s) > 0; {
		top := len(s) - 1
		u := s[top]
		if len(g[u]) > 0 {
			arcs := g[u]
			w := arcs[0]
			s = append(s, w)
			g[u] = arcs[1:]
		} else {
			p[m] = u
			m--
			s = s[:top]
		}
	}
	if m > 0 {
		return nil, errors.New("not Eulerian")
	}
	return p, nil
}

// EulerianPath finds an Eulerian path in a directed multigraph.
//
// * If g has no nodes, result is nil, nil.
//
// * If g has an Eulerian path, result is an Eulerian path with err = nil.
// The path result is a list of nodes, where the first node is start.
//
// * Otherwise, result is nil, error
//
// Internally, EulerianPath copies the entire graph g.
// See EulerianPathD for a more space efficient version.
func (g AdjacencyList) EulerianPath() ([]int, error) {
	ind := g.InDegree()
	var start int
	for n, to := range g {
		if len(to) > ind[n] {
			start = n
			break
		}
	}
	c, m := g.Copy()
	return c.EulerianPathD(m, start)
}

// EulerianPathD finds an Eulerian path in a directed multigraph.
//
// EulerianPathD is destructive on its receiver g.  See EulerianPath for
// a non-destructive version.
//
// Argument m must be the correct size, or number of arcs in g.  Argument
// start must be a valid start node for the path.
//
// * If g has no nodes, result is nil, nil.
//
// * If g has an Eulerian path, result is an Eulerian path with err = nil.
// The path result is a list of nodes, where the first node is start.
//
// * Otherwise, result is nil, error
func (g AdjacencyList) EulerianPathD(m, start int) ([]int, error) {
	if len(g) == 0 {
		return nil, nil
	}
	p := make([]int, m+1)
	for s := []int{start}; len(s) > 0; {
		top := len(s) - 1
		u := s[top]
		if len(g[u]) > 0 {
			arcs := g[u]
			w := arcs[0]
			s = append(s, w)
			g[u] = arcs[1:]
		} else {
			p[m] = u
			m--
			s = s[:top]
		}
	}
	if m > 0 {
		return nil, errors.New("no Eulerian path")
	}
	return p, nil
}
