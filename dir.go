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
	"sort"
)

// DAGMaxLenPath finds a maximum length path in a directed acyclic graph.
//
// Argument ordering must be a topological ordering of g.
func (g AdjacencyList) DAGMaxLenPath(ordering []NI) (path []NI) {
	// dynamic programming. visit nodes in reverse order. for each, compute
	// longest path as one plus longest of 'to' nodes.
	// Visits each arc once.  O(m).
	//
	// Similar code in label.go
	var n NI
	mlp := make([][]NI, len(g)) // index by node number
	for i := len(ordering) - 1; i >= 0; i-- {
		fr := ordering[i] // node number
		to := g[fr]
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
func (g AdjacencyList) EulerianCycle() ([]NI, error) {
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
func (g AdjacencyList) EulerianCycleD(m int) ([]NI, error) {
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
func (g AdjacencyList) EulerianCycleUndirD(m int) ([]NI, error) {
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
func (g AdjacencyList) EulerianPath() ([]NI, error) {
	ind := g.InDegree()
	var start NI
	for n, to := range g {
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
// Argument m must be the correct size, or number of arcs in g.  Argument
// start must be a valid start node for the path.
//
// * If g has no nodes, result is nil, nil.
//
// * If g has an Eulerian path, result is an Eulerian path with err = nil.
// The path result is a list of nodes, where the first node is start.
//
// * Otherwise, result is nil, error
func (g AdjacencyList) EulerianPathD(m int, start NI) ([]NI, error) {
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

// HasParallelSort identifies if a graph contains parallel arcs, multiple arcs
// that lead from a node to the same node.
//
// If the graph has parallel arcs, the results fr and to represent an example
// where there are parallel arcs from node fr to node to.
//
// If there are no parallel arcs, the method returns false -1 -1.
//
// Multiple loops on a node count as parallel arcs.
//
// "Sort" in the method name indicates that sorting is used to detect parallel
// arcs.  Compared to method HasParallelMap, this may give better performance
// for small or sparse graphs but will have asymtotically worse performance for
// large dense graphs.
func (g AdjacencyList) HasParallelSort() (has bool, fr, to NI) {
	var t NodeList
	for n, to := range g {
		if len(to) == 0 {
			continue
		}
		t = append(t[:0], to...)
		sort.Sort(t)
		t0 := t[0]
		for _, to := range t[1:] {
			if to == t0 {
				return true, NI(n), t0
			}
			t0 = to
		}
	}
	return false, -1, -1
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
// of the path will be the same node, indicating a cycle.
//
// Paths are sent on the returned channel.  The channel is closed after all
// paths are sent.
func (g AdjacencyList) MaximalNonBranchingPaths() chan []NI {
	ch := make(chan []NI)
	go g.mnbp(ch)
	return ch
}

func (g AdjacencyList) mnbp(ch chan []NI) {
	ind := g.InDegree()
	var uv big.Int
	OneBits(&uv, len(g))
	for v, vTo := range g {
		if !(ind[v] == 1 && len(vTo) == 1) {
			for _, w := range vTo {
				n := []NI{NI(v), w}
				uv.SetBit(&uv, v, 0)
				uv.SetBit(&uv, int(w), 0)
				wTo := g[w]
				for ind[w] == 1 && len(wTo) == 1 {
					u := wTo[0]
					n = append(n, u)
					uv.SetBit(&uv, int(u), 0)
					w = u
					wTo = g[w]
				}
				ch <- n // path
			}
		}
	}
	for b := uv.BitLen(); b > 0; b = uv.BitLen() {
		v := NI(b - 1)
		n := []NI{v}
		for w := v; ; {
			w = g[w][0]
			uv.SetBit(&uv, int(w), 0)
			n = append(n, w)
			if w == v {
				break
			}
		}
		ch <- n // isolated cycle
	}
	close(ch)
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

// Transpose, for directed graphs, constructs a new adjacency list that is
// the transpose of g.
//
// For every arc from->to of g, the result will have an arc to->from.
// Transpose also counts arcs as it traverses and returns m the number of arcs
// in g (equal to the number of arcs in the result.)
func (g AdjacencyList) Transpose() (t AdjacencyList, m int) {
	t = make(AdjacencyList, len(g))
	for n, nbs := range g {
		for _, nb := range nbs {
			t[nb] = append(t[nb], NI(n))
			m++
		}
	}
	return
}
