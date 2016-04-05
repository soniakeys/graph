// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package df_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/df"
)

func ExampleBits() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {1},
	}
	var b graph.Bits
	fmt.Println("3210")
	fmt.Println("----")
	df.Search(g, 0, df.Bits(&b), df.OkNodeVisitor(func(graph.NI) bool {
		fmt.Printf("%04b\n", &b)
		return true
	}, nil))
	// Output:
	// 3210
	// ----
	// 0001
	// 0011
	// 0111
	// 1111
}

func ExampleVisitor() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {1},
	}
	var ok bool
	df.Search(g, 0, df.OkNodeVisitor(func(n graph.NI) (ok bool) {
		fmt.Println("visit", n)
		return true
	}, &ok))
	fmt.Println(ok)
	// Output:
	// visit 0
	// visit 1
	// visit 2
	// visit 3
	// true
}

func ExampleVisitor_earlyTermination() {
	//   0-->3
	//  / \
	// 1-->2
	g := graph.AdjacencyList{
		0: {1, 2, 3},
		1: {2},
		3: {},
	}
	var ok bool
	df.Search(g, 0, df.OkNodeVisitor(func(n graph.NI) (ok bool) {
		fmt.Println("visit", n)
		return n != 2
	}, &ok))
	fmt.Println(ok)
	// Output:
	// visit 0
	// visit 1
	// visit 2
	// false
}

var k10 graph.Directed

func init() {
	r := rand.New(rand.NewSource(11))
	k10, _ = graph.KroneckerDir(10, 10, r)
}

func TestK10(t *testing.T) {
	var b graph.Bits
	k10.DepthFirst(0, &b, nil)
	r := b.PopCount()
	t.Log("K10 reached =", r)
	if r < 500 {
		t.Fatal(r)
	}
}

func BenchmarkADF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b graph.Bits
		k10.DepthFirst(0, &b, nil)
	}
}

func BenchmarkDFA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var bm graph.Bits
		df.Search(k10.AdjacencyList, 0, df.Bits(&bm))
	}
}
