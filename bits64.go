// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// +build !386

package graph

import "math/big"

const (
	wordSize = 64
	wordExp  = 6 // 2^6 = 64
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
