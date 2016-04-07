// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleNewBits() {
	x := graph.NewBits(3, 5)
	fmt.Println(x.Slice())
	// Output:
	// [3 5]
}

func ExampleBits_And() {
	x := graph.NewBits(3, 5, 6)
	y := graph.NewBits(4, 5, 6)
	x.And(x, y)
	fmt.Println(x.Slice())
	// Output:
	// [5 6]
}

func ExampleBits_AllNot() {
	x := graph.NewBits(3, 5)
	x.AllNot(8, x)
	fmt.Println(x.Slice())
	// Output:
	// [0 1 2 4 6 7]
}

func ExampleBits_AndNot() {
	x := graph.NewBits(0, 2, 3, 5)
	y := graph.NewBits(1, 3)
	x.AndNot(x, y)
	fmt.Println(x.Slice())
	// Output:
	// [0 2 5]
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

func ExampleBits_Clear() {
	x := graph.NewBits(3, 5)
	fmt.Println(x.Slice())
	x.Clear()
	fmt.Println(x.Slice())
	// Output:
	// [3 5]
	// []
}

func ExampleBits_Format() {
	b := graph.NewBits(0, 2, 3)
	fmt.Printf("%4b\n", b)
	// Output:
	// 1101
}

func ExampleBits_From() {
	b := graph.NewBits(0, 2, 128, 129)
	for n := b.From(0); n >= 0; n = b.From(n + 1) {
		fmt.Println(n)
	}
	// Output:
	// 0
	// 2
	// 128
	// 129
}

func ExampleBits_Or() {
	x := graph.NewBits(3, 5, 6)
	y := graph.NewBits(4, 5, 6)
	x.Or(x, y)
	fmt.Println(x.Slice())
	// Output:
	// [3 4 5 6]
}

func ExampleBits_PopCount() {
	b := graph.NewBits(0, 2, 128)
	fmt.Println(b.PopCount())
	// Output:
	// 3
}

func ExampleBits_Set() {
	x := graph.NewBits(0, 2)
	var z graph.Bits
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
	x := graph.NewBits(0, 2)
	y := graph.NewBits(129)
	var z graph.Bits

	fmt.Println(x.PopCount(), "bits, single =", x.Single())
	fmt.Println(y.PopCount(), "bit,  single =", y.Single())
	fmt.Println(z.PopCount(), "bits, single =", z.Single())
	// Output:
	// 2 bits, single = false
	// 1 bit,  single = true
	// 0 bits, single = false
}

func ExampleBits_Slice() {
	b := graph.NewBits(0, 2)
	fmt.Println(b.Slice())
	// Output:
	// [0 2]
}

func ExampleBits_Xor() {
	x := graph.NewBits(3, 5, 6)
	y := graph.NewBits(4, 5, 6)
	x.Xor(x, y)
	fmt.Println(x.Slice())
	// Output:
	// [3 4]
}

func ExampleBits_Zero() {
	x := graph.NewBits(0, 2)
	y := graph.NewBits(129)
	var z graph.Bits

	fmt.Println(x.PopCount(), "bits, zero =", x.Zero())
	fmt.Println(y.PopCount(), "bit,  zero =", y.Zero())
	fmt.Println(z.PopCount(), "bits, zero =", z.Zero())
	// Output:
	// 2 bits, zero = false
	// 1 bit,  zero = false
	// 0 bits, zero = true
}
