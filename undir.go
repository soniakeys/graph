// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// undir.go
//
// Methods specific to undirected graphs.
// Doc for each method should specifically say undirected.

package ed

import (
	"math/big"
)

// ConnectedComponents, for undirected graphs, determines the connected
// components in g.
//
// Returned is a slice with a single representative node from each connected
// component.
func (g AdjacencyList) ConnectedComponents() []int {
	var r []int
	var c big.Int
	var df func(int)
	df = func(n int) {
		c.SetBit(&c, n, 1)
		for _, nb := range g[n] {
			if c.Bit(nb) == 0 {
				df(nb)
			}
		}
	}
	for n := range g {
		if c.Bit(n) == 0 {
			r = append(r, n)
			df(n)
		}
	}
	return r
}
