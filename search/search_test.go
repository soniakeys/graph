// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package search_test

import (
	"fmt"
	"math/rand"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

func ExampleArcVisitor() {
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
	search.DepthFirst(g, 0, search.ArcVisitor(func(n graph.NI, x int) {
		fmt.Println(n, "->", g[n][x])
	}))
	// Output:
	// 0 -> 1
	// 1 -> 2
	// 2 -> 3
	// 3 -> 1
	// 0 -> 2
}

func ExampleNodeVisitor() {
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
	search.DepthFirst(g, 0, search.NodeVisitor(func(n graph.NI) {
		fmt.Println(n)
	}))
	// Output:
	// 0
	// 1
	// 2
	// 3
}

func ExampleOkArcVisitor() {
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
	search.DepthFirst(g, 0, search.OkArcVisitor(func(n graph.NI, x int) bool {
		fmt.Println(n, "->", g[n][x])
		return n < g[n][x]
	}))
	// Output:
	// 0 -> 1
	// 1 -> 2
	// 2 -> 3
	// 3 -> 1
}

func ExampleOkNodeVisitor() {
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
	search.DepthFirst(g, 0, search.OkNodeVisitor(func(n graph.NI) bool {
		fmt.Println(n)
		return n != 2
	}))
	// Output:
	// 0
	// 1
	// 2
}

func ExamplePathBits() {
	//   0
	//  / \
	// 1   2
	// ^   |
	// |   v
	// \---3
	g := graph.AdjacencyList{
		0: {1, 2},
		2: {3},
		3: {1},
	}
	b := bits.New(len(g))
	fmt.Println("node  path bits")
	fmt.Println("      (3210)")
	fmt.Println("----   ----")
	search.DepthFirst(g, 0, search.PathBits(&b), search.NodeVisitor(func(n graph.NI) {
		fmt.Printf("%4d   %s\n", n, &b)
	}))
	// Output:
	// node  path bits
	//       (3210)
	// ----   ----
	//    0   0001
	//    1   0011
	//    2   0101
	//    3   1101
}

func ExampleVisited() {
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
	b := bits.New(len(g))
	fmt.Println("3210")
	fmt.Println("----")
	search.DepthFirst(g, 0, search.Visited(&b),
		search.NodeVisitor(func(graph.NI) {
			fmt.Println(b)
		}))
	// Output:
	// 3210
	// ----
	// 0001
	// 0011
	// 0111
	// 1111
}

func ExampleRand() {
	//         0
	//         |
	// -------------------
	// | | | | | | | | | |
	// 1 2 3 4 5 6 7 8 9 10
	g := graph.AdjacencyList{
		0:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		10: nil,
	}
	search.DepthFirst(g, 0,
		search.Rand(rand.New(rand.NewSource(7))),
		search.NodeVisitor(func(n graph.NI) {
			fmt.Println(n)
		}))
	// Output:
	// 0
	// 3
	// 1
	// 6
	// 4
	// 2
	// 7
	// 10
	// 9
	// 5
	// 8
}
