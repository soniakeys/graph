// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"fmt"
	"math/big"
)

var one = big.NewInt(1)

// trailingZeros returns the number of trailing 0 bits in v.
//
// If v is 0, it returns 0.
func trailingZeros(v big.Word) int {
	return deBruijnBits[v&-v*deBruijnMultiple>>deBruijnShift]
}

type Bits struct {
	i big.Int
}

func (z *Bits) AndNot(x, y *Bits) {
	z.i.AndNot(&x.i, &y.i)
}

func (b Bits) Bit(n NI) uint {
	return b.i.Bit(int(n))
}

func (z *Bits) Clear() {
	z.i.SetInt64(0)
}

func (b Bits) Format(s fmt.State, ch rune) {
	b.i.Format(s, ch)
}

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

// NextOne facilitates iteration over the one bits of a big.Int.
//
// It returns the position of the first one bit at or after position i.
// It returns -1 if there is no one bit at or after position i.
//
// To iterate over the one bits of a big.Int, call with i = 0 to get the
// the first one bit, then call with i+1 to get successive one bits.
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

// PopCount returns the number of one bits.
func (b Bits) PopCount() (c int) {
	for _, w := range b.i.Bits() {
		for w != 0 {
			w &= w - 1
			c++
		}
	}
	return
}

func (z *Bits) Set(x *Bits) {
	z.i.Set(&x.i)
}

// It's a utility function useful for initializing a bitmap of a graph
// to all ones; that is, with a bit set to 1 for each node of the graph.
func (z *Bits) SetAll(n int) {
	z.i.Sub(z.i.Lsh(one, uint(n)), one)
}

func (z *Bits) SetBit(n NI, b uint) {
	z.i.SetBit(&z.i, int(n), b)
}

// like PopCount, but stop as soon as two are found
func (b Bits) Single() bool {
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

func (b Bits) Slice() (s []NI) {
	// (alternative implementation might use Popcount and make to get the
	// exact cap slice up front.  unclear if this would be better.)
	b.Iterate(func(n NI) bool {
		s = append(s, n)
		return true
	})
	return
}

func (b Bits) Zero() bool {
	return len(b.i.Bits()) == 0
}
