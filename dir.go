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
	e := newEulerian(g, m)
	for e.s >= 0 {
		v := e.top() // v is node that starts cycle
		e.push()
		// if Eulerian, we'll always come back to starting node
		if e.top() != v {
			return nil, errors.New("not balanced")
		}
		e.keep()
	}
	if e.uv.BitLen() > 0 {
		return nil, errors.New("not strongly connected")
	}
	return e.p, nil
}

// EulerianCycleUndirD is a bit of an experiment.
//
// It is about the same as EulerianCycleD, but modified for an undirected
// multigraph.
//
// Parameter m in this case must be the size of the undirected graph -- the
// number of edges -- which is the number of arcs / 2.
//
// It works, but contains an extra loop that I think spoils the time
// complexity.  Probably still pretty fast in practice, but a different
// graph representation might be better.
func (g AdjacencyList) EulerianCycleUndirD(m int) ([]int, error) {
	if len(g) == 0 {
		return nil, nil
	}
	e := newEulerian(g, m)
	for e.s >= 0 {
		v := e.top()
		e.pushUndir() // call modified method
		if e.top() != v {
			return nil, errors.New("not balanced")
		}
		e.keep()
	}
	if e.uv.BitLen() > 0 {
		return nil, errors.New("not strongly connected")
	}
	return e.p, nil
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
	e := newEulerian(g, m)
	e.p[0] = start
	// unlike EulerianCycle, the first path doesn't have be a cycle.
	e.push()
	e.keep()
	for e.s >= 0 {
		start = e.top()
		e.push()
		// paths after the first must be cycles though
		// (as long as there are nodes on the stack)
		if e.top() != start {
			return nil, errors.New("no Eulerian path")
		}
		e.keep()
	}
	if e.uv.BitLen() > 0 {
		return nil, errors.New("no Eulerian path")
	}
	return e.p, nil
}

// starting at the node on the top of the stack, follow arcs until stuck.
// mark nodes visited, push nodes on stack, remove arcs from g.
func (e *eulerian) push() {
	for u := e.top(); ; {
		e.uv.SetBit(&e.uv, u, 0) // reset unvisited bit
		arcs := e.g[u]
		if len(arcs) == 0 {
			return // stuck
		}
		w := arcs[0] // follow first arc
		e.s++        // push followed node on stack
		e.p[e.s] = w
		e.g[u] = arcs[1:] // consume arc
		u = w
	}
}

// like push, but for for undirected graphs.
func (e *eulerian) pushUndir() {
	for u := e.top(); ; {
		e.uv.SetBit(&e.uv, u, 0)
		arcs := e.g[u]
		if len(arcs) == 0 {
			return
		}
		w := arcs[0]
		e.s++
		e.p[e.s] = w
		e.g[u] = arcs[1:] // consume arc
		// here is the only difference, consume reciprocal arc as well:
		a2 := e.g[w]
		for x, rx := range a2 {
			if rx == u { // here it is
				last := len(a2) - 1
				a2[x] = a2[last]   // someone else gets the seat
				e.g[w] = a2[:last] // and it's gone.
				break
			}
		}
		u = w
	}
}

// starting with the node on top of the stack, move nodes with no arcs.
func (e *eulerian) keep() {
	for e.s >= 0 {
		n := e.top()
		if len(e.g[n]) > 0 {
			break
		}
		e.p[e.m] = n
		e.s--
		e.m--
	}
}

type eulerian struct {
	g  AdjacencyList // working copy of graph, it gets consumed
	m  int           // number of arcs in g, updated as g is consumed
	uv big.Int       // unvisited
	// low end of p is stack of unfinished nodes
	// high end is finished path
	p []int // stack + path
	s int   // stack pointer
}

func (e *eulerian) top() int {
	return e.p[e.s]
}

func newEulerian(g AdjacencyList, m int) *eulerian {
	e := &eulerian{
		g: g,
		m: m,
		p: make([]int, m+1),
	}
	e.uv.Lsh(one, uint(len(g)))
	e.uv.Sub(&e.uv, one)
	return e
}

func (g AdjacencyList) MaximalNonBranchingPaths() (p [][]int) {
	ind := g.InDegree()
	var uv big.Int
	uv.Lsh(one, uint(len(g)))
	uv.Sub(&uv, one)
	for v, vTo := range g {
		if !(ind[v] == 1 && len(vTo) == 1) {
			for _, w := range vTo {
				n := []int{v, w}
				uv.SetBit(&uv, v, 0)
				uv.SetBit(&uv, w, 0)
				wTo := g[w]
				for ind[w] == 1 && len(wTo) == 1 {
					u := wTo[0]
					n = append(n, u)
					uv.SetBit(&uv, u, 0)
					w = u
					wTo = g[w]
				}
				// path
				p = append(p, n)
			}
		}
	}
	for b := uv.BitLen(); b > 0; b = uv.BitLen() {
		v := b - 1
		n := []int{v}
		for w := v; ; {
			w = g[w][0]
			uv.SetBit(&uv, w, 0)
			n = append(n, w)
			if w == v {
				break
			}
		}
		// isolated cycle
		p = append(p, n)
	}
	return p
}
