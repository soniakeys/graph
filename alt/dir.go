// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt

import (
	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

// TarjanCycles emits all elementary cycles of a directed graph.
//
// See graph.Cycles for Johnson's algorithm.
func TarjanCycles(g graph.Directed, emit func([]graph.NI) bool) {
	// Implementation of "Enumeration of the elementary circuits of a directed
	// graph", by Robert Tarjan, TR 72-145, Cornell University, September 1972.
	a, _ := g.AdjacencyList.Copy()
	mark := bits.New(len(a))
	var point, marked []graph.NI
	var s graph.NI
	var backtrack func(graph.NI, string) (bool, bool)
	backtrack = func(v graph.NI, indent string) (f, ok bool) {
		point = append(point, v)
		mark.SetBit(int(v), 1)
		marked = append(marked, v)
		to := a[v]
		for i := 0; i < len(to); {
			w := to[i]
			if w < s {
				// "delete w from A(v)"
				last := len(to) - 1
				to[i] = to[last]
				to = to[:last]
				a[v] = to
			} else {
				i++
				if w == s {
					if !emit(point) {
						return f, false
					}
					f = true
				} else if mark.Bit(int(w)) == 0 {
					switch g, ok := backtrack(w, indent+"  "); {
					case !ok:
						return f, false
					case g:
						f = true
					}
				}
			}
		}
		// "f=true if an elementary circuit continuing the partial path
		// on the stack has been found"
		if f {
			top := len(marked) - 1
			for ; marked[top] != v; top-- {
				u := marked[top]
				marked = marked[:top]
				mark.SetBit(int(u), 0)
			}
			marked = marked[:top]
			mark.SetBit(int(v), 0)
		}
		point = point[:len(point)-1]
		return f, true
	}
	for s = graph.NI(0); int(s) < len(a); s++ {
		if _, ok := backtrack(graph.NI(s), ""); !ok {
			return
		}
		for top := len(marked) - 1; top >= 0; top-- {
			mark.SetBit(int(marked[top]), 0)
			marked = marked[:top]
		}
	}
}

// TarjanCycles emits all elementary cycles of a directed graph.
//
// See graph.Cycles for Johnson's algorithm.
func TarjanCyclesLabeled(g graph.LabeledDirected, emit func([]graph.Half) bool) {
	a, _ := g.LabeledAdjacencyList.Copy()
	mark := bits.New(len(a))
	var half []graph.Half
	var marked []graph.NI
	var s graph.NI
	var backtrack func(graph.NI, string) (bool, bool)
	backtrack = func(v graph.NI, indent string) (f, ok bool) {
		mark.SetBit(int(v), 1)
		marked = append(marked, v)
		to := a[v]
		for i := 0; i < len(to); {
			w := to[i]
			if w.To < s {
				last := len(to) - 1
				to[i] = to[last]
				to = to[:last]
				a[v] = to
			} else {
				i++
				if w.To == s {
					if !emit(append(half, w)) {
						return f, false
					}
					f = true
				} else if mark.Bit(int(w.To)) == 0 {
					half = append(half, w)
					switch g, ok := backtrack(w.To, indent+"  "); {
					case !ok:
						return f, false
					case g:
						f = true
					}
					half = half[:len(half)-1]
				}
			}
		}
		if f {
			top := len(marked) - 1
			for ; marked[top] != v; top-- {
				u := marked[top]
				marked = marked[:top]
				mark.SetBit(int(u), 0)
			}
			marked = marked[:top]
			mark.SetBit(int(v), 0)
		}
		return f, true
	}
	for s = graph.NI(0); int(s) < len(a); s++ {
		if _, ok := backtrack(graph.NI(s), ""); !ok {
			return
		}
		for top := len(marked) - 1; top >= 0; top-- {
			mark.SetBit(int(marked[top]), 0)
			marked = marked[:top]
		}
	}
}

// Note well:  The backing slice for the node list passed to emit is reused
// across emit calls.  If you need to retain the node list you must copy it.
func SCCKosaraju(g graph.Directed, emit func([]graph.NI) bool) {
	// 1. For each vertex u of the graph, mark u as unvisited. Let L be empty.
	a := g.AdjacencyList
	t := make(graph.AdjacencyList, len(a)) // transpose graph
	vis := make([]bool, len(a))
	L := make([]graph.NI, len(a))
	x := len(L) // index for filling L in reverse order

	// 2. recursive subroutine:
	var Visit func(graph.NI)
	Visit = func(u graph.NI) {
		if !vis[u] {
			vis[u] = true
			for _, v := range a[u] {
				Visit(v)
				t[v] = append(t[v], u) // construct transpose
			}
			x--
			L[x] = u
		}
	}
	// 2. For each vertex u of the graph do Visit(u)
	for u := range a {
		Visit(graph.NI(u))
	}
	var c []graph.NI // result, the component assignment
	// 3: recursive subroutine:
	var Assign func(graph.NI)
	Assign = func(u graph.NI) {
		if vis[u] { // repurpose vis to mean "unassigned"
			vis[u] = false
			c = append(c, u)
			for _, v := range t[u] {
				Assign(v)
			}
		}
	}
	// 3: For each element u of L in order, do Assign(u)
	for _, u := range L {
		Assign(u)
		if len(c) > 0 {
			if !emit(c) {
				return
			}
		}
		c = c[:0] // reuse slice
	}
}

// SCCPathBased identifies strongly connected components in a directed graph
// using a path-based algorithm.
//
// The method calls the emit argument for each component identified.  Each
// component is a list of nodes.  The emit function must return true for the
// method to continue identifying components.  If emit returns false, the
// method returns immediately.
//
// Note well:  The backing slice for the node list passed to emit is reused
// across emit calls.  If you need to retain the node list you must copy it.
func SCCPathBased(g graph.Directed, emit func([]graph.NI) bool) {
	a := g.AdjacencyList
	var S []graph.NI
	var B []int
	I := make([]int, len(a))
	for i := range I {
		I[i] = -1
	}
	c := len(a)
	var scc []graph.NI
	var df func(graph.NI)
	df = func(v graph.NI) {
		I[v] = len(S)
		B = append(B, len(S))
		S = append(S, v)
		for _, w := range a[v] {
			if I[w] < 0 {
				df(w)
			} else {
				for last := len(B) - 1; I[w] < B[last]; last-- {
					B = B[:last]
				}
			}
		}
		if last := len(B) - 1; I[v] == B[last] {
			scc = scc[:0]
			B = B[:last]
			for I[v] <= len(S) {
				last := len(S) - 1
				I[S[last]] = c
				scc = append(scc, S[last])
				S = S[:last]
			}
			c++
			emit(scc)
		}
	}
	for v := range a {
		if I[v] < 0 {
			df(graph.NI(v))
		}
	}
}

// SCCTarjan identifies strongly connected components in a directed graph using
// Tarjan's algorithm.
//
// The method calls the emit argument for each component identified.  Each
// component is a list of nodes.  The emit function must return true for the
// method to continue identifying components.  If emit returns false, the
// method returns immediately.
//
// A property of the algorithm is that components are emitted in reverse
// topological order of the condensation.
// (See https://en.wikipedia.org/wiki/Strongly_connected_component#Definitions
// for description of condensation.)
//
// Note well:  The backing slice for the node list passed to emit is reused
// across emit calls.  If you need to retain the node list you must copy it.
func SCCTarjan(g graph.Directed, emit func([]graph.NI) bool) {
	// See "Depth-first search and linear graph algorithms", Robert Tarjan,
	// SIAM J. Comput. Vol. 1, No. 2, June 1972.
	//
	// Implementation here from Wikipedia pseudocode.
	a := g.AdjacencyList
	indexed := bits.New(len(a))
	stacked := bits.New(len(a))
	index := make([]int, len(a))
	lowlink := make([]int, len(a))
	x := 0
	var S, c []graph.NI
	var sc func(graph.NI) bool
	sc = func(n graph.NI) bool {
		index[n] = x
		indexed.SetBit(int(n), 1)
		lowlink[n] = x
		x++
		S = append(S, n)
		stacked.SetBit(int(n), 1)
		for _, nb := range a[n] {
			if indexed.Bit(int(nb)) == 0 {
				if !sc(nb) {
					return false
				}
				if lowlink[nb] < lowlink[n] {
					lowlink[n] = lowlink[nb]
				}
			} else if stacked.Bit(int(nb)) == 1 {
				if index[nb] < lowlink[n] {
					lowlink[n] = index[nb]
				}
			}
		}
		if lowlink[n] == index[n] {
			c = c[:0]
			for {
				last := len(S) - 1
				w := S[last]
				S = S[:last]
				stacked.SetBit(int(w), 0)
				c = append(c, w)
				if w == n {
					if !emit(c) {
						return false
					}
					break
				}
			}
		}
		return true
	}
	for n := range a {
		if indexed.Bit(n) == 0 && !sc(graph.NI(n)) {
			return
		}
	}
}
