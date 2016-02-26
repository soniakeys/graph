package graph

import (
	"math/big"
	"math/rand"
	"time"
)

// KroneckerDir generates a Kronecker-like random directed graph.
//
// The returned graph g is simple and has no isolated nodes but is not
// necessarily fully connected.  The number of of nodes will be <= 2^scale,
// and will be near 2^scale for typical values of arcFactor, >= 2.
// ArcFactor * 2^scale arcs are generated, although loops and duplicate arcs
// are rejected.
//
// Return value m is the number of arcs retained in the result graph.
func KroneckerDir(scale uint, arcFactor float64) (g AdjacencyList, m int) {
	return kronecker(scale, arcFactor, true)
}

// KroneckerUndir generates a Kronecker-like random undirected graph.
//
// The returned graph g is simple and has no isolated nodes but is not
// necessarily fully connected.  The number of of nodes will be <= 2^scale,
// and will be near 2^scale for typical values of edgeFactor, >= 2.
// EdgeFactor * 2^scale edges are generated, although loops and duplicate edges
// are rejected.
//
// Return value m is the number of arcs--not edges--retained in the result
// graph.
func KroneckerUndir(scale uint, edgeFactor float64) (g Undirected, m int) {
	al, as := kronecker(scale, edgeFactor, false)
	return Undirected{al}, as
}

// Styled after the Graph500 example code.  Not well tested currently.
// Graph500 example generates undirected only.  No idea if the directed variant
// here is meaningful or not.
func kronecker(scale uint, edgeFactor float64, dir bool) (g AdjacencyList, m int) {
	rand.Seed(time.Now().Unix())
	N := 1 << scale                      // node extent
	M := int(edgeFactor*float64(N) + .5) // number of arcs/edges to generate
	a, b, c := 0.57, 0.19, 0.19          // initiator probabilities
	ab := a + b
	cNorm := c / (1 - ab)
	aNorm := a / ab
	ij := make([][2]int, M)
	var bm big.Int
	var nNodes int
	for k := range ij {
		var i, j int
		for b := 1; b < N; b <<= 1 {
			if rand.Float64() > ab {
				i |= b
				if rand.Float64() > cNorm {
					j |= b
				}
			} else if rand.Float64() > aNorm {
				j |= b
			}
		}
		if bm.Bit(i) == 0 {
			bm.SetBit(&bm, i, 1)
			nNodes++
		}
		if bm.Bit(j) == 0 {
			bm.SetBit(&bm, j, 1)
			nNodes++
		}
		r := rand.Intn(k + 1) // shuffle edges as they are generated
		ij[k] = ij[r]
		ij[r] = [2]int{i, j}
	}
	p := rand.Perm(nNodes) // mapping to shuffle IDs of non-isolated nodes
	px := 0
	r := make([]NI, N)
	for i := range r {
		if bm.Bit(i) == 1 {
			r[i] = NI(p[px]) // fill lookup table
			px++
		}
	}
	g = make(AdjacencyList, nNodes)
ij:
	for _, e := range ij {
		if e[0] == e[1] {
			continue // skip loops
		}
		ri, rj := r[e[0]], r[e[1]]
		for _, nb := range g[ri] {
			if nb == rj {
				continue ij // skip parallel edges
			}
		}
		g[ri] = append(g[ri], rj)
		m++
		if !dir {
			g[rj] = append(g[rj], ri)
			m++
		}
	}
	return
}
