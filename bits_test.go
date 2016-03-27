// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleBits_AndNot() {
	var x, y graph.Bits
	x.SetBit(0, 1)
	x.SetBit(2, 1)
	x.SetBit(128, 1)
	x.SetBit(129, 1)

	y.SetBit(128, 1)
	x.AndNot(x, y)
	fmt.Println(x.Slice())
	// Output:
	// [0 2 129]
}

func ExampleBits_Bit() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	for n := graph.NI(0); n < 4; n++ {
		fmt.Println("bit", n, "=", b.Bit(n))
	}
	// Output:
	// bit 0 = 1
	// bit 1 = 0
	// bit 2 = 1
	// bit 3 = 0
}

func ExampleBits_Format() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	b.SetBit(3, 1)
	fmt.Printf("%4b\n", b)
	// Output:
	// 1101
}

func ExampleBits_Iterate() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	b.SetBit(128, 1)
	b.SetBit(129, 1)
	b.Iterate(func(n graph.NI) bool {
		fmt.Println(n)
		return true
	})
	// Output:
	// 0
	// 2
	// 128
	// 129
}

func ExampleBits_NextOne() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	b.SetBit(128, 1)
	b.SetBit(129, 1)
	for n := b.NextOne(0); n >= 0; n = b.NextOne(n + 1) {
		fmt.Println(n)
	}
	// Output:
	// 0
	// 2
	// 128
	// 129
}

func ExampleBits_PopCount() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	b.SetBit(128, 1)
	fmt.Println(b.PopCount())
	// Output:
	// 3
}

func ExampleBits_Set() {
	var x, z graph.Bits
	x.SetBit(0, 1)
	x.SetBit(2, 1)
	z.Set(x)
	fmt.Println(z.Slice())
	// Output:
	// [0 2]
}

func ExampleBits_SetAll() {
	g := make(graph.AdjacencyList, 5)
	var b graph.Bits
	b.SetAll(len(g))
	fmt.Printf("%b\n", b)
	// Output:
	// 11111
}

func ExampleBits_SetBit() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	fmt.Println(b.Slice())
	// Output:
	// [0 2]
}

func ExampleBits_Single() {
	var x, y, z graph.Bits
	x.SetBit(0, 1)
	x.SetBit(2, 1)

	y.SetBit(129, 1)

	fmt.Println(x.PopCount(), "bits, single =", x.Single())
	fmt.Println(y.PopCount(), "bit,  single =", y.Single())
	fmt.Println(z.PopCount(), "bits, single =", z.Single())
	// Output:
	// 2 bits, single = false
	// 1 bit,  single = true
	// 0 bits, single = false
}

func ExampleBits_Slice() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	fmt.Println(b.Slice())
	// Output:
	// [0 2]
}

func ExampleBits_Zero() {
	var x, y, z graph.Bits
	x.SetBit(0, 1)
	x.SetBit(2, 1)

	y.SetBit(129, 1)

	fmt.Println(x.PopCount(), "bits, zero =", x.Zero())
	fmt.Println(y.PopCount(), "bit,  zero =", y.Zero())
	fmt.Println(z.PopCount(), "bits, zero =", z.Zero())
	// Output:
	// 2 bits, zero = false
	// 1 bit,  zero = false
	// 0 bits, zero = true
}
