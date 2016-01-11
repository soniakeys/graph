// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_Bipartite() {
	g := graph.AdjacencyList{
		0: {3},
		1: {3},
		2: {3, 4},
		3: {0, 1, 2},
		4: {2},
	}
	b, c1, c2, oc := g.Bipartite(0)
	if b {
		fmt.Println(
			strconv.FormatInt(c1.Int64(), 2),
			strconv.FormatInt(c2.Int64(), 2))
	}
	g[3] = append(g[3], 4)
	g[4] = append(g[4], 3)
	b, c1, c2, oc = g.Bipartite(0)
	if !b {
		fmt.Println(oc)
	}
	// Output:
	// 111 11000
	// [3 4 2]
}

func ExampleAdjacencyList_BronKerbosch1() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.AdjacencyList
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	for c := range g.BronKerbosch1() {
		fmt.Println(c)
	}
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
	var g graph.AdjacencyList
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	for c := range g.BronKerbosch2(g.BKPivotMaxDegree) {
		fmt.Println(c)
	}
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
	var g graph.AdjacencyList
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	for c := range g.BronKerbosch3(g.BKPivotMaxDegree) {
		fmt.Println(c)
	}
	// Output:
	// [0 4]
	// [3 4]
	// [4 5]
	// [2 3]
	// [1 2 5]
}

func ExampleAdjacencyList_UndirectedCopy_simple() {
	//    0
	//   / \
	//  1   2
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {},
		2: {},
	}
	fmt.Println(g.IsUndirected())
	u := g.UndirectedCopy()
	for fr, to := range u {
		fmt.Println(fr, to)
	}
	ok, _, _ := u.IsUndirected()
	fmt.Println(ok)
	// Output:
	// false 0 1
	// 0 [1 2]
	// 1 [0]
	// 2 [0]
	// true
}

func ExampleAdjacencyList_UndirectedCopy_loopMultigraph() {
	//  0--\   1--->2
	//  |  |   |<---|
	//  \--/   |<---|
	g := graph.AdjacencyList{
		0: {0},
		1: {2},
		2: {1, 1},
	}
	fmt.Println(g.IsUndirected())
	u := g.UndirectedCopy()
	for fr, to := range u {
		fmt.Println(fr, to)
	}
	ok, _, _ := u.IsUndirected()
	fmt.Println(ok)
	// Output:
	// false 2 1
	// 0 [0]
	// 1 [2 2]
	// 2 [1 1]
	// true
}

func ExampleAdjacencyList_IsUndirected() {
	//    0  2--\
	//   /   |  |
	//  1    \--/
	g := graph.AdjacencyList{
		0: {1},
		1: {0},
		2: {2},
	}
	ud, _, _ := g.IsUndirected()
	fmt.Println(ud)
	//   0
	//  /
	// 1<--2
	g = graph.AdjacencyList{
		0: {1},
		1: {0},
		2: {1},
	}
	fmt.Println(g.IsUndirected())
	// Output:
	// true
	// false 2 1
}

func ExampleAdjacencyList_IsTreeUndirected() {
	//  1--\
	//  |  |
	//  \--/   0   5
	//        /   / \
	//       2   3---4
	g := graph.AdjacencyList{
		1: {1},
		0: {2},
		2: {0},
		5: {3, 4},
		3: {4, 5},
		4: {3, 5},
	}
	fmt.Println(g.IsTreeUndirected(1))
	fmt.Println(g.IsTreeUndirected(2))
	fmt.Println(g.IsTreeUndirected(3))
	// Output:
	// false
	// true
	// false
}

func ExampleAdjacencyList_ConnectedComponentReps() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	g := graph.AdjacencyList{
		0: {3, 4},
		1: {5},
		3: {0, 4},
		4: {0, 3},
		5: {1},
	}
	fmt.Println(g.ConnectedComponentReps())
	// Output:
	// [0 1 2] [3 2 1]
}

func ExampleAdjacencyList_ConnectedComponentReps_collectingBits() {
	g := graph.AdjacencyList{
		0: {3, 4},
		1: {5},
		3: {0, 4},
		4: {0, 3},
		5: {1},
	}
	rep, order := g.ConnectedComponentReps()
	fmt.Println("543210  rep  order")
	fmt.Println("------  ---  -----")
	for i, r := range rep {
		var bits big.Int
		g.DepthFirst(r, &bits, nil)
		fmt.Printf("%0*b   %d     %d\n", len(g), &bits, r, order[i])
	}
	// Output:
	// 543210  rep  order
	// ------  ---  -----
	// 011001   0     3
	// 100010   1     2
	// 000100   2     1
}

func ExampleAdjacencyList_ConnectedComponentReps_collectingLists() {
	g := graph.AdjacencyList{
		0: {3, 4},
		1: {5},
		3: {0, 4},
		4: {0, 3},
		5: {1},
	}
	rep, _ := g.ConnectedComponentReps()
	for _, r := range rep {
		var m []int
		g.DepthFirst(r, nil, func(n int) bool {
			m = append(m, n)
			return true
		})
		fmt.Println(m)
	}
	// Output:
	// [0 3 4]
	// [1 5]
	// [2]
}

func ExampleAdjacencyList_ConnectedComponentBits() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	g := graph.AdjacencyList{
		0: {3, 4},
		1: {5},
		3: {0, 4},
		4: {0, 3},
		5: {1},
	}
	f := g.ConnectedComponentBits()
	fmt.Println("o  543210")
	fmt.Println("-  ------")
	for o, b := f(); o > 0; o, b = f() {
		fmt.Printf("%d  %0*b\n", o, len(g), &b)
	}
	// Output:
	// o  543210
	// -  ------
	// 3  011001
	// 2  100010
	// 1  000100
}

func ExampleAdjacencyList_ConnectedComponentLists() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	g := graph.AdjacencyList{
		0: {3, 4},
		1: {5},
		3: {0, 4},
		4: {0, 3},
		5: {1},
	}
	f := g.ConnectedComponentLists()
	for l := f(); l != nil; l = f() {
		fmt.Println(l)
	}
	// Output:
	// [0 3 4]
	// [1 5]
	// [2]
}

func ExampleTarjanBiconnectedComponents() {
	g := graph.AdjacencyList{}
	g.AddEdge(3, 4)
	g.AddEdge(3, 2)
	g.AddEdge(2, 4)
	g.AddEdge(2, 5)
	g.AddEdge(2, 1)
	g.AddEdge(5, 1)
	g.AddEdge(6, 1)
	g.AddEdge(6, 5)
	g.AddEdge(7, 1)
	g.AddEdge(7, 9)
	g.AddEdge(7, 8)
	g.AddEdge(9, 8)
	for _, bcc := range g.TarjanBiconnectedComponents() {
		fmt.Println("Edges:")
		for _, e := range bcc {
			fmt.Println(e)
		}
	}
	// Output:
	// Edges:
	// {4 2}
	// {3 4}
	// {2 3}
	// Edges:
	// {6 1}
	// {5 6}
	// {5 1}
	// {2 5}
	// {1 2}
	// Edges:
	// {8 7}
	// {9 8}
	// {7 9}
	// Edges:
	// {1 7}
}

/* shelved
func ExampleBiconnectedComponents_Find() {
	g := graph.AdjacencyList{
		0:  {1, 7},
		1:  {2, 4, 0},
		2:  {3, 1},
		3:  {2, 4},
		4:  {3, 1},
		5:  {6, 12},
		6:  {5, 12, 8},
		7:  {8, 0},
		8:  {6, 7, 9, 10},
		9:  {8},
		10: {8, 13, 11},
		11: {10},
		12: {5, 6, 13},
		13: {12, 10},
	}
	b := graph.NewBiconnectedComponents(g)
	b.Find(0)
	fmt.Println("n: cut from")
	for n, f := range b.From {
		fmt.Printf("%d: %d %d\n",
			n, b.Cuts.Bit(n), f)
	}
	fmt.Println("Leaves:", b.Leaves)
	// Output:
	// n: cut from
	// 0: 1 -1
	// 1: 1 0
	// 2: 0 1
	// 3: 0 2
	// 4: 0 3
	// 5: 0 6
	// 6: 0 8
	// 7: 1 0
	// 8: 1 7
	// 9: 0 8
	// 10: 1 13
	// 11: 0 10
	// 12: 0 5
	// 13: 0 12
	// Leaves: [4 11 9]
}
*/

func ExampleAdjacencyList_Degeneracy() {
	//
	//   /---2
	//  0--./|\
	//  |\ /\| \
	//  | .  6  3--5
	//  |/ \/| /
	//  1--'\|/
	//   \---4
	//
	// Same graph redrawn to show ordering:
	//
	//          /-----\
	//         /  /--\ \
	//      /-+--+-\  \|
	//  5--3--4--6--2--1--0
	//         \  \  \---/|\
	//          \  \-----/ |
	//           \--------/
	//
	g := graph.AdjacencyList{
		0: {1, 2, 4, 6},
		1: {0, 2, 4, 6},
		2: {0, 1, 3, 6},
		3: {2, 4, 5},
		4: {0, 1, 3, 6},
		5: {3},
		6: {0, 1, 2, 4},
	}
	k, ord, cores := g.Degeneracy()
	fmt.Println("Degeneracy:", k)
	fmt.Println("Ordering:", ord)
	fmt.Println("0-core:", ord[:cores[0]])
	for k := 1; k < len(cores); k++ {
		fmt.Printf("%d-core: %d\n", k, ord[cores[k-1]:cores[k]])
	}
	// Output:
	// Degeneracy: 3
	// Ordering: [5 3 4 6 2 1 0]
	// 0-core: []
	// 1-core: [5]
	// 2-core: [3]
	// 3-core: [4 6 2 1 0]
}
