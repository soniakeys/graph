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

// Directed represents a directed graph.
//
// Directed methods generally rely on the graph being directed, specifically
// that arcs do not have reciprocals.
type Directed struct {
	AdjacencyList // embedded to include AdjacencyList methods
}

// DAGMaxLenPath finds a maximum length path in a directed acyclic graph.
//
// Argument ordering must be a topological ordering of g.
func (g Directed) DAGMaxLenPath(ordering []NI) (path []NI) {
	// dynamic programming. visit nodes in reverse order. for each, compute
	// longest path as one plus longest of 'to' nodes.
	// Visits each arc once.  O(m).
	//
	// Similar code in label.go
	var n NI
	mlp := make([][]NI, len(g.AdjacencyList)) // index by node number
	for i := len(ordering) - 1; i >= 0; i-- {
		fr := ordering[i] // node number
		to := g.AdjacencyList[fr]
		if len(to) == 0 {
			continue
		}
		mt := to[0]
		for _, to := range to[1:] {
			if len(mlp[to]) > len(mlp[mt]) {
				mt = to
			}
		}
		p := append([]NI{mt}, mlp[mt]...)
		mlp[fr] = p
		if len(p) > len(path) {
			n = fr
			path = p
		}
	}
	return append([]NI{n}, path...)
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
func (g Directed) EulerianCycle() ([]NI, error) {
	c, m := g.Copy()
	return c.EulerianCycleD(m)
}

// EulerianCycleD finds an Eulerian cycle in a directed multigraph.
//
// EulerianCycleD is destructive on its receiver g.  See EulerianCycle for
// a non-destructive version.
//
// Argument ma must be the correct arc size, or number of arcs in g.
//
// * If g has no nodes, result is nil, nil.
//
// * If g is Eulerian, result is an Eulerian cycle with err = nil.
// The cycle result is a list of nodes, where the first and last
// nodes are the same.
//
// * Otherwise, result is nil, error
func (g Directed) EulerianCycleD(ma int) ([]NI, error) {
	if len(g.AdjacencyList) == 0 {
		return nil, nil
	}
	e := newEulerian(g.AdjacencyList, ma)
	for e.s >= 0 {
		v := e.top() // v is node that starts cycle
		e.push()
		// if Eulerian, we'll always come back to starting node
		if e.top() != v {
			return nil, errors.New("not balanced")
		}
		e.keep()
	}
	if len(e.uv.Bits()) > 0 {
		return nil, errors.New("not strongly connected")
	}
	return e.p, nil
}

// EulerianCycleD for undirected graphs is a bit of an experiment.
//
// It is about the same as the directed version, but modified for an undirected
// multigraph.
//
// Parameter m in this case must be the size of the undirected graph -- the
// number of edges.  Use Undirected.Size if the size is unknown.
//
// It works, but contains an extra loop that I think spoils the time
// complexity.  Probably still pretty fast in practice, but a different
// graph representation might be better.
func (g Undirected) EulerianCycleD(m int) ([]NI, error) {
	if len(g.AdjacencyList) == 0 {
		return nil, nil
	}
	e := newEulerian(g.AdjacencyList, m)
	for e.s >= 0 {
		v := e.top()
		e.pushUndir() // call modified method
		if e.top() != v {
			return nil, errors.New("not balanced")
		}
		e.keep()
	}
	if len(e.uv.Bits()) > 0 {
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
func (g Directed) EulerianPath() ([]NI, error) {
	ind := g.InDegree()
	var start NI
	for n, to := range g.AdjacencyList {
		if len(to) > ind[n] {
			start = NI(n)
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
// Argument ma must be the correct arc size, or number of arcs in g.
// Argument start must be a valid start node for the path.
//
// * If g has no nodes, result is nil, nil.
//
// * If g has an Eulerian path, result is an Eulerian path with err = nil.
// The path result is a list of nodes, where the first node is start.
//
// * Otherwise, result is nil, error
func (g Directed) EulerianPathD(ma int, start NI) ([]NI, error) {
	if len(g.AdjacencyList) == 0 {
		return nil, nil
	}
	e := newEulerian(g.AdjacencyList, ma)
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
	if len(e.uv.Bits()) > 0 {
		return nil, errors.New("no Eulerian path")
	}
	return e.p, nil
}

// starting at the node on the top of the stack, follow arcs until stuck.
// mark nodes visited, push nodes on stack, remove arcs from g.
func (e *eulerian) push() {
	for u := e.top(); ; {
		e.uv.SetBit(&e.uv, int(u), 0) // reset unvisited bit
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
		e.uv.SetBit(&e.uv, int(u), 0)
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
	p []NI // stack + path
	s int  // stack pointer
}

func (e *eulerian) top() NI {
	return e.p[e.s]
}

func newEulerian(g AdjacencyList, m int) *eulerian {
	e := &eulerian{
		g: g,
		m: m,
		p: make([]NI, m+1),
	}
	OneBits(&e.uv, len(g))
	return e
}

// MaximalNonBranchingPaths finds all paths in a directed graph that are
// "maximal" and "non-branching".
//
// A non-branching path is one where path nodes other than the first and last
// have exactly one arc leading to the node and one arc leading from the node,
// thus there is no possibility to branch away to a different path.
//
// A maximal non-branching path cannot be extended to a longer non-branching
// path by including another node at either end.
//
// In the case of a cyclic non-branching path, the first and last elements
// of the path will be the same node, indicating an isolated cycle.
//
// The method calls the emit argument for each path or isolated cycle in g,
// as long as emit returns true.  If emit returns false,
// MaximalNonBranchingPaths returns immediately.
func (g Directed) MaximalNonBranchingPaths(emit func([]NI) bool) {
	ind := g.InDegree()
	var uv big.Int
	OneBits(&uv, len(g.AdjacencyList))
	for v, vTo := range g.AdjacencyList {
		if !(ind[v] == 1 && len(vTo) == 1) {
			for _, w := range vTo {
				n := []NI{NI(v), w}
				uv.SetBit(&uv, v, 0)
				uv.SetBit(&uv, int(w), 0)
				wTo := g.AdjacencyList[w]
				for ind[w] == 1 && len(wTo) == 1 {
					u := wTo[0]
					n = append(n, u)
					uv.SetBit(&uv, int(u), 0)
					w = u
					wTo = g.AdjacencyList[w]
				}
				if !emit(n) { // n is a path
					return
				}
			}
		}
	}
	for b := NextOne(&uv, 0); b >= 0; b = NextOne(&uv, b+1) {
		v := NI(b)
		n := []NI{v}
		for w := v; ; {
			w = g.AdjacencyList[w][0]
			uv.SetBit(&uv, int(w), 0)
			n = append(n, w)
			if w == v {
				break
			}
		}
		if !emit(n) { // n is an isolated cycle
			return
		}
	}
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

// Transpose constructs a new adjacency list with all arcs reversed.
//
// For every arc from->to of g, the result will have an arc to->from.
// Transpose also counts arcs as it traverses and returns ma the number of arcs
// in g (equal to the number of arcs in the result.)
func (g Directed) Transpose() (t Directed, ma int) {
	ta := make(AdjacencyList, len(g.AdjacencyList))
	for n, nbs := range g.AdjacencyList {
		for _, nb := range nbs {
			ta[nb] = append(ta[nb], NI(n))
			ma++
		}
	}
	return Directed{ta}, ma
}
