// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_ArcSize() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.AdjacencyList{
		2: {0, 1},
	}
	fmt.Println(g.ArcSize()) // simple graph

	// with reciprocals now.
	u := g.UndirectedCopy()
	// common term "size" for undirected graph is number of undirected edges.
	// size, m = ArcSize() / 2 here, but only because there are no loops.
	fmt.Println(u.ArcSize())

	g[1] = []graph.NI{1} // add a loop
	//   2
	//  / \
	// 0   1---\
	//      \--/
	fmt.Println(g.ArcSize())

	// loops have no reciprocals.  ArcSize() / 2 no longer meaningful.
	fmt.Println(g.UndirectedCopy().ArcSize())
	// Output:
	// 2
	// 4
	// 3
	// 5
}

func ExampleAdjacencyList_Balanced() {
	// 2
	// |
	// v
	// 0----->1
	g := graph.AdjacencyList{
		2: {0},
		0: {1},
	}
	fmt.Println(g.Balanced())

	// 0<--\
	// |    \
	// v     \
	// 1----->2
	g[1] = []graph.NI{2}
	fmt.Println(g.Balanced())
	// Output:
	// false
	// true
}

func ExampleAdjacencyList_Cyclic() {
	//   0
	//  / \
	// 1-->2-->3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {},
	}
	cyclic, _, _ := g.Cyclic()
	fmt.Println(cyclic)

	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g[3] = []graph.NI{1}
	fmt.Println(g.Cyclic())

	// Output:
	// false
	// true 3 1
}

func ExampleAdjacencyList_DepthFirst() {
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
	ok := g.DepthFirst(0, nil, func(n graph.NI) (ok bool) {
		fmt.Println("visit", n)
		return true
	})
	fmt.Println(ok)
	// Output:
	// visit 0
	// visit 1
	// visit 2
	// visit 3
	// true
}

func ExampleAdjacencyList_DepthFirst_earlyTermination() {
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
	ok := g.DepthFirst(0, nil, func(n graph.NI) bool {
		fmt.Println("visit", n)
		return n != 2
	})
	fmt.Println(ok)
	// Output:
	// visit 0
	// visit 1
	// visit 2
	// false
}

func ExampleAdjacencyList_DepthFirst_bitmap() {
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
	var vis big.Int
	fmt.Println("3210")
	fmt.Println("----")
	g.DepthFirst(0, &vis, func(graph.NI) bool {
		fmt.Printf("%04b\n", &vis)
		return true
	})
	// Output:
	// 3210
	// ----
	// 0001
	// 0011
	// 0111
	// 1111
}

func TestAdjacencyList_DepthFirst_bothNil(t *testing.T) {
	// for coverage
	var g graph.AdjacencyList
	if g.DepthFirst(0, nil, nil) {
		t.Fatal("DepthFirst both nil must return false")
	}
}

func ExampleAdjacencyList_HasLoop() {
	g := graph.AdjacencyList{
		1: {0},
	}
	fmt.Println(g.HasLoop()) // -1 result means no loop
	g[0] = []graph.NI{0}     // add loop
	fmt.Println(g.HasLoop()) // 0 result means loop on node 0
	// Output:
	// -1
	// 0
}

func ExampleAdjacencyList_Tarjan() {
	// /---0---\
	// |   |\--/
	// |   v
	// |   5<=>4---\
	// |   |   |   |
	// v   v   |   |
	// 7<=>6   |   |
	//     |   v   v
	//     \-->3<--2
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
	}
	for c := range g.Tarjan() {
		fmt.Println(c)
	}
	// Output:
	// [1 3 2]
	// [7 6]
	// [4 5]
	// [0]
}

func ExampleAdjacencyList_TarjanForward() {
	// /---0---\
	// |   |\--/
	// |   v
	// |   5<=>4---\
	// |   |   |   |
	// v   v   |   |
	// 7<=>6   |   |
	//     |   v   v
	//     \-->3<--2
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
	}
	for _, c := range g.TarjanForward() {
		fmt.Println(c)
	}
	// Output:
	// [0]
	// [4 5]
	// [7 6]
	// [1 3 2]
}
