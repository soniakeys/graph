// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt

import (
	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

// SCCPathBased identifies strongly connected components in a directed graph
// using a path-based algorithm.
//
// The method calls the emit argument for each component identified.  Each
// component is a list of nodes.  The emit function must return true for the
// method to continue identifying components.  If emit returns false, the
// method returns immediately.
func SCCPathBased(g graph.Directed, emit func([]graph.NI) bool) {
	a := g.AdjacencyList
	var S []graph.NI
	var B []int
	I := make([]int, len(a))
	for i := range I {
		I[i] = -1
	}
	c := len(a)
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
			var scc []graph.NI
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
	var S []graph.NI
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
			var c []graph.NI
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
