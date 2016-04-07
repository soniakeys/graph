// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// undir_cg_test.go -- tests for code in undir_cg.go.
//
// These are tests on the labeled versions of methods.
//
// See also undir_ro_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleLabeledUndirected_Bipartite() {
	// 0 1 2
	//  \|/|
	//   3 4
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	b, c1, c2, _ := g.Bipartite(0)
	if b {
		fmt.Println("n:  43210")
		fmt.Printf("c1: %05b\n", c1)
		fmt.Printf("c2: %05b\n", c2)
	}
	// Output:
	// n:  43210
	// c1: 00111
	// c2: 11000
}

func ExampleLabeledUndirected_Bipartite_oddCycle() {
	// 0 1  2
	//  \|/ |
	//   3--4
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	b, _, _, oc := g.Bipartite(0)
	if !b {
		fmt.Println("odd cycle:", oc)
	}
	// Output:
	// odd cycle: [3 4 2]
}

func ExampleLabeledUndirected_BronKerbosch1() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
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

func ExampleLabeledUndirected_BKPivotMaxDegree() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
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

func ExampleLabeledUndirected_BKPivotMinP() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
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

func ExampleLabeledUndirected_BronKerbosch2() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
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

func ExampleLabeledUndirected_BronKerbosch3() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
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

func ExampleLabeledUndirected_ConnectedComponentBits() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	f := g.ConnectedComponentBits()
	fmt.Println("o  543210")
	fmt.Println("-  ------")
	for o, b := f(); o > 0; o, b = f() {
		fmt.Printf("%d  %0*b\n", o, len(g.LabeledAdjacencyList), &b)
	}
	// Output:
	// o  543210
	// -  ------
	// 3  011001
	// 2  100010
	// 1  000100
}

func ExampleLabeledUndirected_ConnectedComponentLists() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	f := g.ConnectedComponentLists()
	for l := f(); l != nil; l = f() {
		fmt.Println(l)
	}
	// Output:
	// [0 3 4]
	// [1 5]
	// [2]
}

func ExampleLabeledUndirected_ConnectedComponentReps() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	fmt.Println(g.ConnectedComponentReps())
	// Output:
	// [0 1 2] [3 2 1]
}

func ExampleLabeledUndirected_ConnectedComponentReps_collectingBits() {
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	rep, order := g.ConnectedComponentReps()
	fmt.Println("543210  rep  order")
	fmt.Println("------  ---  -----")
	for i, r := range rep {
		var bits graph.Bits
		g.DepthFirst(r, &bits, nil)
		fmt.Printf("%0*b   %d     %d\n",
			len(g.LabeledAdjacencyList), &bits, r, order[i])
	}
	// Output:
	// 543210  rep  order
	// ------  ---  -----
	// 011001   0     3
	// 100010   1     2
	// 000100   2     1
}

func ExampleLabeledUndirected_ConnectedComponentReps_collectingLists() {
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	rep, _ := g.ConnectedComponentReps()
	for _, r := range rep {
		var m []graph.NI
		g.DepthFirst(r, nil, func(n graph.NI) bool {
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

func ExampleLabeledUndirected_Degeneracy() {
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
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{0, 6}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{1, 4}, 0)
	g.AddEdge(graph.Edge{1, 6}, 0)
	g.AddEdge(graph.Edge{6, 2}, 0)
	g.AddEdge(graph.Edge{6, 4}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{3, 5}, 0)
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

func ExampleLabeledUndirected_Degree() {
	// 0---1--\
	//      \-/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{1, 1}, 0)
	fmt.Println(g.Degree(0))
	fmt.Println(g.Degree(1))
	// Output:
	// 1
	// 3
}

func ExampleLabeledUndirected_FromList() {
	//    4   3
	//   / \
	//  2   1
	//       \
	//        0
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{4, 1}, 0)
	g.AddEdge(graph.Edge{1, 0}, 0)
	f, cycle := g.FromList(4)
	if cycle >= 0 {
		return
	}
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %3d\n", n, e.From, e.Len)
	}
	fmt.Println("MaxLen:  ", f.MaxLen)
	// Output:
	// n  from  path len
	// 0    1    3
	// 1    4    2
	// 2    4    2
	// 3   -1    0
	// 4   -1    1
	// MaxLen:   3
}

func ExampleLabeledUndirected_FromList_cycle() {
	//    0
	//   / \
	//  1   2
	//     / \
	//    3---4
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	_, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: 2
}

func ExampleLabeledUndirected_FromList_loop() {
	//    0
	//   / \ /-\
	//  1   2--/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{2, 2}, 0)
	_, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: 2
}

func ExampleLabeledUndirected_FromList_loopDisconnected() {
	//    0
	//   /   /-\
	//  1   2--/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{2, 2}, 0)
	f, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: -1
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %3d\n", n, e.From, e.Len)
	}
	// Output:
	// cycle: -1
	// n  from  path len
	// 0   -1    1
	// 1    0    2
	// 2   -1    0
}

func ExampleLabeledUndirected_FromList_multigraph() {
	//    0
	//   / \
	//  1   2==3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	_, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: 3
}

func ExampleLabeledUndirected_IsConnected() {
	// undirected graph:
	//   0
	//  / \
	// 1   2
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// connected:  true
}

func ExampleLabeledUndirected_IsConnected_notConnected() {
	// undirected graph:
	//   0   1
	//  / \
	// 2   3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{0, 3}, 0)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// connected:  false
}
func ExampleLabeledUndirected_IsTree() {
	//  0--\
	//  |  |
	//  \--/   1   3
	//        /   / \
	//       2   4---5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 0}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{3, 5}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	fmt.Println(g.IsTree(0))
	fmt.Println(g.IsTree(1))
	fmt.Println(g.IsTree(3))
	// Output:
	// false false
	// true false
	// false false
}

func ExampleLabeledUndirected_Size() {
	//   0--\
	//  / \-/
	// 1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 0}, 0)
	g.AddEdge(graph.Edge{0, 1}, 0)
	fmt.Println("Size:", g.Size())
	fmt.Printf("(Arc size = %d)\n", g.ArcSize())
	// Output:
	// Size: 2
	// (Arc size = 3)
}
