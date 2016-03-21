// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// +build 386

package graph

import "math/big"

const (
	wordSize = 32
	wordExp  = 5 // 2^5 = 32
)

// NextOne facilitates iteration over the one bits of a big.Int.
//
// It returns the position of the first one bit at or after position i.
// It returns -1 if there is no one bit at or after position i.
//
// To iterate over the one bits of a big.Int, call with i = 0 to get the
// the first one bit, then call with i+1 to get successive one bits.
func NextOne(b *big.Int, i int) int {
	words := b.Bits()
	x := i >> wordExp // x now index of word containing bit i.
	if x >= len(words) {
		return -1
	}
	// test for 1 in this word at or after i
	if wx := words[x] >> (uint(i) & (wordSize - 1)); wx != 0 {
		return i + trailingZeros(wx)
	}
	x++
	for y, wy := range words[x:] {
		if wy != 0 {
			return (x+y)<<wordExp | trailingZeros(wy)
		}
	}
	return -1
}

// reference: http://graphics.stanford.edu/~seander/bithacks.html
const deBruijn32Multiple = 0x077CB531
const deBruijn32Shift = 27

var deBruijn32Bits = []int{
	0, 1, 28, 2, 29, 14, 24, 3, 30, 22, 20, 15, 25, 17, 4, 8,
	31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9,
}

// trailingZeros returns the number of trailing 0 bits in v.
//
// If v is 0, it returns 0.
func trailingZeros(v big.Word) int {
	return deBruijn32Bits[v&-v*deBruijn32Multiple>>deBruijn32Shift]
}
