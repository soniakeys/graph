// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"fmt"
	"math/big"
)

// Bits is bitmap, or bitset, intended to store a single bit of information
// per node of a graph.
//
// The current implementation is backed by a big.Int and so is a reference
// type in the same way a big.Int is.
type Bits struct {
	i big.Int
}

// AndNot sets z = x &^ y and returns z.
func (z *Bits) AndNot(x, y *Bits) {
	z.i.AndNot(&x.i, &y.i)
}

// Bit returns the value of the n'th bit of x.
func (b Bits) Bit(n NI) uint {
	return b.i.Bit(int(n))
}

// Clear sets all bits to 0.
func (z *Bits) Clear() {
	z.i.SetInt64(0)
}

// Format satisfies fmt.Formatter for fmt.Printf and related methods.
//
// graph.Bits format exactly like big.Ints.
func (b Bits) Format(s fmt.State, ch rune) {
	b.i.Format(s, ch)
}

// Iterate calls Visitor v for each bit with a value of 1, in order
// from lowest bit to highest bit.
//
// Iteration continues to the highest bit as long as v returns true.
// It stops if v returns false.
//
// Iterate returns true normally.  It returns false if v returns false.
//
// Bit values should not be modified during iteration, by the visitor function
// for example.  See NextOne for an iteration method that allows modification.
func (b Bits) Iterate(v Visitor) bool {
	for x, w := range b.i.Bits() {
		if w != 0 {
			t := trailingZeros(w)
			i := t // index in w of next 1 bit
			for {
				if !v(NI(x<<wordExp | i)) {
					return false
				}
				w >>= uint(t + 1)
				if w == 0 {
					break
				}
				t = trailingZeros(w)
				i += 1 + t
			}
		}
	}
	return true
}

// NextOne facilitates iteration over 1 bits.
//
// It returns the position of the first one bit at or after position n.
// It returns -1 if there is no one bit at or after position n.
//
// To iterate over the one bits, call with n = 0 to get the the first
// one bit, then call with the result + 1 to get successive one bits.
// Unlike the Iterate method, this technique is stateless and so allows
// bits to be changed between successive calls.
func (b Bits) NextOne(n NI) NI {
	words := b.i.Bits()
	i := int(n)
	x := i >> wordExp // x now index of word containing bit i.
	if x >= len(words) {
		return -1
	}
	// test for 1 in this word at or after n
	if wx := words[x] >> (uint(i) & (wordSize - 1)); wx != 0 {
		return n + NI(trailingZeros(wx))
	}
	x++
	for y, wy := range words[x:] {
		if wy != 0 {
			return NI((x+y)<<wordExp | trailingZeros(wy))
		}
	}
	return -1
}

// PopCount returns the number of 1 bits.
func (b Bits) PopCount() (c int) {
	// algorithm selected to be efficient for sparse bit sets.
	for _, w := range b.i.Bits() {
		for w != 0 {
			w &= w - 1
			c++
		}
	}
	return
}

// Set sets all bits of z to the bits of x.
func (z *Bits) Set(x *Bits) {
	z.i.Set(&x.i)
}

var one = big.NewInt(1)

// SetAll sets z to have n 1 bits.
//
// It's useful for initializing z to have a 1 for each node of a graph.
func (z *Bits) SetAll(n int) {
	z.i.Sub(z.i.Lsh(one, uint(n)), one)
}

// SetBit sets the n'th bit to b, where be is a 0 or 1.
func (z *Bits) SetBit(n NI, b uint) {
	z.i.SetBit(&z.i, int(n), b)
}

// Single returns true if b has exactly one 1 bit.
func (b Bits) Single() bool {
	// like PopCount, but stop as soon as two are found
	c := 0
	for _, w := range b.i.Bits() {
		for w != 0 {
			w &= w - 1
			c++
			if c == 2 {
				return false
			}
		}
	}
	return c == 1
}

// Slice returns a slice with the positons of each 1 bit.
func (b Bits) Slice() (s []NI) {
	// (alternative implementation might use Popcount and make to get the
	// exact cap slice up front.  unclear if that would be better.)
	b.Iterate(func(n NI) bool {
		s = append(s, n)
		return true
	})
	return
}

// Zero returns true if there are no 1 bits.
func (b Bits) Zero() bool {
	return len(b.i.Bits()) == 0
}

// trailingZeros returns the number of trailing 0 bits in v.
//
// If v is 0, it returns 0.
func trailingZeros(v big.Word) int {
	return deBruijnBits[v&-v*deBruijnMultiple>>deBruijnShift]
}
