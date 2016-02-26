// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/big"
	"math/rand"
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

func ExampleAdjacencyList_ArcSize_handshakingLemma() {
	// undirected graph with three edges:
	//   0
	//   |
	//   1   2
	//  / \
	// 3   4
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(1, 3)
	g.AddEdge(1, 4)
	// for undirected g without loops, degree == out-degree == len(to)
	degSum := 0
	for _, to := range g.AdjacencyList {
		degSum += len(to)
	}
	// for undirected g without loops, ArcSize is 2 * number of edges
	fmt.Println(degSum, "==", g.ArcSize())
	// Output:
	// 6 == 6
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

func ExampleAdjacencyList_BoundsOk() {
	var g graph.AdjacencyList
	ok, _, _ := g.BoundsOk() // zero value adjacency list is valid
	fmt.Println(ok)
	g = graph.AdjacencyList{
		0: {9},
	}
	fmt.Println(g.BoundsOk()) // arc 0 to 9 invalid with only one node
	// Output:
	// true
	// false 0 9
}

func ExampleAdjacencyList_BronKerbosch1() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch1(func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleAdjacencyList_BKPivotMaxDegree() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleAdjacencyList_BKPivotMinP() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch2(g.BKPivotMinP, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleAdjacencyList_BronKerbosch2() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleAdjacencyList_BronKerbosch3() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch3(g.BKPivotMaxDegree, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [3 4]
	// [4 5]
	// [2 3]
	// [1 2 5]
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

func ExampleAdjacencyList_DepthFirstRandom() {
	//     ----0-----
	//    /    |     \
	//   1     2      3
	//  /|\   /|\   / | \
	// 4 5 6 7 8 9 10 11 12
	g := graph.AdjacencyList{
		0:  {1, 2, 3},
		1:  {4, 5, 6},
		2:  {7, 8, 9},
		3:  {10, 11, 12},
		12: {},
	}
	r := rand.New(rand.NewSource(12))
	f := func(n graph.NI) (ok bool) {
		fmt.Println("visit", n)
		return true
	}
	g.DepthFirstRandom(0, nil, f, r)
	// Output:
	// visit 0
	// visit 1
	// visit 6
	// visit 4
	// visit 5
	// visit 3
	// visit 12
	// visit 11
	// visit 10
	// visit 2
	// visit 9
	// visit 7
	// visit 8
}

func ExampleAdjacencyList_HasArc() {
	g := graph.AdjacencyList{
		2: {0, 2, 0, 1, 1},
	}
	fmt.Println(g.HasArc(2, 1))
	// Output:
	// true 3
}

func ExampleAdjacencyList_HasLoop_loop() {
	g := graph.AdjacencyList{
		2: {2},
	}
	fmt.Println(g.HasLoop())
	// Output:
	// true 2
}

func ExampleAdjacencyList_HasLoop_noLoop() {
	g := graph.AdjacencyList{
		1: {0},
	}
	lp, _ := g.HasLoop()
	fmt.Println("has loop:", lp)
	// Output:
	// has loop: false
}

func ExampleAdjacencyList_HasParallelMap_parallelArcs() {
	g := graph.AdjacencyList{
		1: {0, 0},
	}
	// result true 1 0 means parallel arcs from node 1 to node 0
	fmt.Println(g.HasParallelMap())
	// Output:
	// true 1 0
}

func ExampleAdjacencyList_HasParallelMap_noParallelArcs() {
	g := graph.AdjacencyList{
		1: {0},
	}
	fmt.Println(g.HasParallelMap()) // result false -1 -1 means no parallel arc
	// Output:
	// false -1 -1
}

func ExampleAdjacencyList_InDegree() {
	// arcs directed down:
	//  0     2
	//  |
	//  1
	//  |\
	//  | \
	//  3  4<-\
	//     \--/
	g := graph.AdjacencyList{
		0: {1},
		1: {3, 4},
		4: {4},
	}
	fmt.Println("node:    0 1 2 3 4")
	fmt.Println("in-deg:", g.InDegree())
	// Output:
	// node:    0 1 2 3 4
	// in-deg: [0 1 0 1 2]
}

func ExampleAdjacencyList_IsConnected() {
	// undirected graph:
	//   0
	//  / \
	// 1   2
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	u, _, _ := g.IsUndirected()
	fmt.Println("undirected:", u)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// undirected: true
	// connected:  true
}

func ExampleAdjacencyList_IsConnected_notConnected() {
	// undirected graph:
	//   0   1
	//  / \
	// 2   3
	var g graph.Undirected
	g.AddEdge(0, 2)
	g.AddEdge(0, 3)
	u, _, _ := g.IsUndirected()
	fmt.Println("undirected:", u)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// undirected: true
	// connected:  false
}

func ExampleAdjacencyList_IsSimple() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.AdjacencyList{
		2: {0, 1},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// true -1
}

func ExampleAdjacencyList_IsSimple_loop() {
	// arcs directed down
	//   2
	//  / \
	// 0   1---\
	//      \--/
	g := graph.AdjacencyList{
		2: {0, 1},
		1: {1}, // loop
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 1
}

func ExampleAdjacencyList_IsSimple_parallelArc() {
	// arcs directed down
	//   2
	//  /|\
	//  |/ \
	//  0   1
	g := graph.AdjacencyList{
		2: {0, 1, 0},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 2
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
	g.Tarjan(func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
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

func ExampleAdjacencyList_UndirectedDegree() {
	// 0---1--\
	//      \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(1, 1)
	fmt.Println(g.Degree(0))
	fmt.Println(g.Degree(1))
	// Output:
	// 1
	// 3
}
