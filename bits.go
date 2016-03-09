// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"math/big"
	"unsafe"
)

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

func init() {
	var w big.Word
	if unsafe.Sizeof(w) != 8 {
		panic("NextOne hard coded to 8 byte words")
	}
}

// NextOne facilitates iteration over the one bits of a big.Int.
//
// It returns the position of the first one bit at or after position i.
// It returns -1 if there is no one bit at or after position i.
//
// To iterate over the one bits of a big.Int, call with i = 0 to get the
// the first one bit, then call with i+1 to get successive one bits.
func NextOne(b *big.Int, i int) int {
	words := b.Bits()
	x := i >> 6 // 8 bytes = 2^6 bits.  x now index of word containing bit i.
	if x >= len(words) {
		return -1
	}
	// a 1 in this word at or after i?
	if wx := words[x] >> uint(i) & 63; wx != 0 {
		return i + trailingZeros(wx)
	}
	x++
	for y, wy := range words[x:] {
		if wy != 0 {
			return (x+y)<<6 | trailingZeros(wy)
		}
	}
	return -1
}

// reference: http://graphics.stanford.edu/~seander/bithacks.html
const deBruijn64Multiple = 0x03f79d71b4ca8b09
const deBruijn64Shift = 58

var deBruijn64Bits = []int{
	0, 1, 56, 2, 57, 49, 28, 3, 61, 58, 42, 50, 38, 29, 17, 4,
	62, 47, 59, 36, 45, 43, 51, 22, 53, 39, 33, 30, 24, 18, 12, 5,
	63, 55, 48, 27, 60, 41, 37, 16, 46, 35, 44, 21, 52, 32, 23, 11,
	54, 26, 40, 15, 34, 20, 31, 10, 25, 14, 19, 9, 13, 8, 7, 6,
}

// trailingZeros returns the number of trailing 0 bits in v.
//
// If v is 0, it returns 0.
func trailingZeros(v big.Word) int {
	return deBruijn64Bits[v&-v*deBruijn64Multiple>>deBruijn64Shift]
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
