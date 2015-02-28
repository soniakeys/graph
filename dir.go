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
	// See "Depth-first search and linear graph algorithms", Robert Tarjan,
	// SIAM J. Comput. Vol. 1, No. 2, June 1972.
	//
	// Implementation here from Wikipedia pseudocode,
	// http://en.wikipedia.org/w/index.php?title=Tarjan%27s_strongly_connected_components_algorithm&direction=prev&oldid=647184742
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

// StronglyConnectedComponents identifies strongly connected components
// in a directed graph.
//
// Algorithm by David J. Pearce, from "An Improved Algorithm for Finding the
// Strongly Connected Components of a Directed Graph".  It is algorithm 3,
// PEA_FIND_SCC2 in
// http://homepages.mcs.vuw.ac.nz/~djp/files/P05.pdf, accessed 22 Feb 2015.
//
// Returned is a list of components, each component is a list of nodes.
/*
func (g AdjacencyList) StronglyConnectedComponents() []int {
	rindex := make([]int, len(g))
	S := []int{}
	index := 1
	c := len(g) - 1
	visit := func(v int) {
		root := true
		rindex[v] = index
		index++
		for _, w := range g[v] {
			if rindex[w] == 0 {
				visit(w)
			}
			if rindex[w] < rindex[v] {
				rindex[v] = rindex[w]
				root = false
			}
		}
		if root {
			index--
			for top := len(S) - 1; top >= 0 && rindex[v] <= rindex[top]; top-- {
				w = rindex[top]
				S = S[:top]
				rindex[w] = c
				index--
			}
			rindex[v] = c
			c--
		} else {
			S = append(S, v)
		}
	}
	for v := range g {
		if rindex[v] == 0 {
			visit(v)
		}
	}
	return rindex
}
*/

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
	uv := new(big.Int).Lsh(one, uint(len(g))) // unvisited
	uv.Sub(uv, one)
	// low end of p is stack of unfinished nodes
	// high end is finished path
	p := make([]int, m+1) // stack + path
	for s := 0; s >= 0; {
		v := p[s] // v is node that starts cycle
		s = g.pushEulerian(uv, p, s)
		if p[s] != v { // if Eulerian, we'll always come back to starting node
			return nil, errors.New("not balanced")
		}
		s, m = g.keepEulerian(p, s, m)
	}
	if uv.BitLen() > 0 {
		return nil, errors.New("not strongly connected")
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
	uv := new(big.Int).Lsh(one, uint(len(g)))
	uv.Sub(uv, one)
	p := make([]int, m+1)
	p[0] = start
	// unlike EulerianCycle, the first path doesn't have be a cycle.
	s := g.pushEulerian(uv, p, 0)
	s, m = g.keepEulerian(p, s, m)
	for s >= 0 {
		start = p[s]
		s = g.pushEulerian(uv, p, s)
		// paths after the first must be cycles though
		if p[s] != start {
			return nil, errors.New("no Eulerian path")
		}
		s, m = g.keepEulerian(p, s, m)
	}
	if uv.BitLen() > 0 {
		return nil, errors.New("no Eulerian path")
	}
	return p, nil
}

// starting at the node on the top of the stack, follow arcs until stuck.
// mark nodes visited, push nodes on stack, remove arcs from g.
// returns new stack pointer.
func (g AdjacencyList) pushEulerian(uv *big.Int, p []int, s int) int {
	for u := p[s]; ; {
		uv.SetBit(uv, u, 0) // reset unvisited bit
		arcs := g[u]
		if len(arcs) == 0 {
			return s // stuck, return stack pointer
		}
		w := arcs[0] // follow first arc
		s++          // push followed node on stack
		p[s] = w
		g[u] = arcs[1:] // consume arc
		u = w
	}
}

// starting with the node on top of the stack, move nodes with no arcs.
// returns new stack pointer and new m.
func (g AdjacencyList) keepEulerian(p []int, s, m int) (int, int) {
	for s >= 0 {
		n := p[s]
		if len(g[n]) > 0 {
			break
		}
		p[m] = n
		s--
		m--
	}
	return s, m
}
