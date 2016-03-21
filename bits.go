// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import "math/big"

var one = big.NewInt(1)

// OneBits sets a big.Int to a number that is all 1s in binary.
//
// It's a utility function useful for initializing a bitmap of a graph
// to all ones; that is, with a bit set to 1 for each node of the graph.
//
// OneBits modifies b, then returns b for convenience.
func OneBits(b *big.Int, n int) *big.Int {
	return b.Sub(b.Lsh(one, uint(n)), one)
}

// PopCount returns the number of one bits in a big.Int.
func PopCount(b *big.Int) (c int) {
	for _, w := range b.Bits() {
		for w != 0 {
			w &= w - 1
			c++
		}
	}
	return
}
