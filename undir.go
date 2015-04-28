// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// undir.go
//
// Methods specific to undirected graphs.
// Doc for each method should specifically say undirected.

package graph

import (
"fmt"
	"math/big"
)

// Edge is an undirected edge between nodes n1 and n2.
type Edge struct{ n1, n2 int }

func (p *AdjacencyList) AddEdge(n1, n2 int) {
	// determine max of the two end points
	max := n1
	if n2 > max {
		max = n2
	}
	// expand graph if needed, to include both
	g := *p
	if max >= len(g) {
		*p = make(AdjacencyList, max+1)
		copy(*p, g)
		g = *p
	}
	// create one half-arc,
	g[n1] = append(g[n1], n2)
	// and except for loops, create the reciprocal
	if n1 != n2 {
		g[n2] = append(g[n2], n1)
	}
}

// ConnectedComponents, for undirected graphs, determines the connected
// components in g.
//
// Returned is a slice with a single representative node from each connected
// component and a parallel slice with the number of nodes in the corresponding
// component.
func (g AdjacencyList) ConnectedComponents() (rep, nNodes []int) {
	var c big.Int
	var nc int
	var df func(int)
	df = func(n int) {
		c.SetBit(&c, n, 1)
		nc++
		for _, nb := range g[n] {
			if c.Bit(nb) == 0 {
				df(nb)
			}
		}
		return
	}
	for n := range g {
		if c.Bit(n) == 0 {
			rep = append(rep, n)
			nc = 0
			df(n)
			nNodes = append(nNodes, nc)
		}
	}
	return
}

// TarjanBiconnectedComponents, for undirected simple graphs.
func (g AdjacencyList) TarjanBiconnectedComponents() (components [][]Edge) {
	// Implemented closely to pseudocode in "Depth-first search and linear
	// graph algorithms", Robert Tarjan, SIAM J. Comput. Vol. 1, No. 2,
	// June 1972.
	//
	// Note Tarjan's "adjacency structure" is graph.AdjacencyList,
	// His "adjacency list" is an element of a graph.AdjacencyList, also
	// termed a "to-list", "neighbor list", or "child list."
	number := make([]int, len(g))
	lowpt := make([]int, len(g))
	var stack []Edge
	var i int
	var biconnect func(int, int)
	biconnect = func(v, u int) {
		i++
		number[v] = i
		lowpt[v] = i
		for _, w := range g[v] {
			if number[w] == 0 {
				stack = append(stack, Edge{v, w})
				biconnect(w, v)
				if lowpt[w] < lowpt[v] {
					lowpt[v] = lowpt[w]
				}
				if lowpt[w] >= number[v] {
					var bcc []Edge
					top := len(stack) - 1
					for number[stack[top].n1] >= number[w] {
						bcc = append(bcc, stack[top])
						stack = stack[:top]
						top--
					}
					bcc = append(bcc, stack[top])
					stack = stack[:top]
					top--
					components = append(components, bcc)
				}
			} else if number[w] < number[v] && w != u {
				stack = append(stack, Edge{v, w})
				if number[w] < lowpt[v] {
					lowpt[v] = number[w]
				}
			}
		}
	}
	for w := range g {
		if number[w] == 0 {
			biconnect(w, 0)
		}
	}
	return
}

/* half-baked.  Read the 72 paper.  Maybe revisit at some point.
type BiconnectedComponents struct {
	Graph  AdjacencyList
	Start  int
	Cuts   big.Int // bitmap of node cuts
	From   []int   // from-tree
	Leaves []int   // leaves of from-tree
}

func NewBiconnectedComponents(g AdjacencyList) *BiconnectedComponents {
	return &BiconnectedComponents{
		Graph: g,
		From:  make([]int, len(g)),
	}
}

func (b *BiconnectedComponents) Find(start int) {
	g := b.Graph
	depth := make([]int, len(g))
	low := make([]int, len(g))
	// reset from any previous run
	b.Cuts.SetInt64(0)
	bf := b.From
	for n := range bf {
		bf[n] = -1
	}
	b.Leaves = b.Leaves[:0]
	d := 1 // depth. d > 0 means visited
	depth[start] = d
	low[start] = d
	d++
	var df func(int, int)
	df = func(from, n int) {
		bf[n] = from
		depth[n] = d
		dn := d
		l := d
		d++
		cut := false
		leaf := true
		for _, nb := range g[n] {
			if depth[nb] == 0 {
				leaf = false
				df(n, nb)
				if low[nb] < l {
					l = low[nb]
				}
				if low[nb] >= dn {
					cut = true
				}
			} else if nb != from && depth[nb] < l {
				l = depth[nb]
			}
		}
		low[n] = l
		if cut {
			b.Cuts.SetBit(&b.Cuts, n, 1)
		}
		if leaf {
			b.Leaves = append(b.Leaves, n)
		}
		d--
	}
	nbs := g[start]
	if len(nbs) == 0 {
		return
	}
	df(start, nbs[0])
	var rc uint
	for _, nb := range nbs[1:] {
		if depth[nb] == 0 {
			rc = 1
			df(start, nb)
		}
	}
	b.Cuts.SetBit(&b.Cuts, start, rc)
	return
}
*/

func (g AdjacencyList) Degeneracy() (k int, ord []int, cores []int) {
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
			fmt.Println("i =", i, "ox =", ox, "ord:", ord[:ox])
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
			if L.Bit(nb) == 1 {
				continue
			}
			dn := d[nb]  // old number of neighbors of nb
			Ddn := D[dn] // nb is in this list
			// remove it from the list
			for wx, w := range Ddn {
				if w == nb {
					last := len(Ddn) - 1
					Ddn[wx], Ddn[last] = Ddn[last], Ddn[wx]
					D[dn] = Ddn[:last]
				}
			}
			dn-- // new number of neighbors
			d[nb] = dn
			// re--add it to it's new list
			D[dn] = append(D[dn], nb)
		}
	}
	cores[k] = len(ord)
	return
}
