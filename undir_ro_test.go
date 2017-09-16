// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// undir_ro_test.go -- tests on undir_RO.go
//
// These are tests on the code-generated unlabeled versions of methods.
//
// If testing prompts changes in the tested method, be sure to edit
// undir_cg.go, go generate to generate the undir_RO.go, and then retest.
// Do not edit undir_RO.go
//
// See also undir_cg_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

func ExampleUndirected_Bipartite() {
	// 0 1 2
	//  \|/|
	//   3 4
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	g.AddEdge(2, 4)
	b, c1, c2, _ := g.Bipartite(0)
	if b {
		fmt.Println("n:  43210")
		fmt.Println("c1:", c1)
		fmt.Println("c2:", c2)
	}
	// Output:
	// n:  43210
	// c1: 00111
	// c2: 11000
}

func ExampleUndirected_Bipartite_oddCycle() {
	// 0 1  2
	//  \|/ |
	//   3--4
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	b, _, _, oc := g.Bipartite(0)
	if !b {
		fmt.Println("odd cycle:", oc)
	}
	// Output:
	// odd cycle: [3 4 2]
}

func ExampleUndirected_BronKerbosch1() {
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
	g.BronKerbosch1(func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleUndirected_BKPivotMaxDegree() {
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
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleUndirected_BKPivotMinP() {
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
	g.BronKerbosch2(g.BKPivotMinP, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleUndirected_BronKerbosch2() {
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
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleUndirected_BronKerbosch3() {
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
	g.BronKerbosch3(g.BKPivotMaxDegree, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [1 2 5]
	// [4 5]
	// [2 3]
	// [3 4]
	// [0 4]
}

func ExampleUndirected_ConnectedComponentBits() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(0, 4)
	g.AddEdge(3, 4)
	g.AddEdge(1, 5)
	f := g.ConnectedComponentBits()
	fmt.Println("o  ma  543210")
	fmt.Println("-  --  ------")
	for {
		n, ma, b := f()
		if n == 0 {
			break
		}
		fmt.Printf("%d  %2d  %s\n", n, ma, b)
	}
	// Output:
	// o  ma  543210
	// -  --  ------
	// 3   6  011001
	// 2   2  100010
	// 1   0  000100
}

func ExampleUndirected_ConnectedComponentLists() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(0, 4)
	g.AddEdge(3, 4)
	g.AddEdge(1, 5)
	f := g.ConnectedComponentLists()
	for {
		l, ma := f()
		if l == nil {
			break
		}
		fmt.Println(l, ma)
	}
	// Output:
	// [0 3 4] 6
	// [1 5] 2
	// [2] 0
}

func ExampleUndirected_ConnectedComponentReps() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(0, 4)
	g.AddEdge(3, 4)
	g.AddEdge(1, 5)
	reps, orders, arcSizes := g.ConnectedComponentReps()
	fmt.Println("reps:    ", reps)
	fmt.Println("orders:  ", orders)
	fmt.Println("arcSizes:", arcSizes)
	// Output:
	// reps:     [0 1 2]
	// orders:   [3 2 1]
	// arcSizes: [6 2 0]
}

func ExampleUndirected_Degeneracy() {
	//   1   ----5
	//  / \ /   / \
	// 0---2---4  |
	//      \   \ /
	//   3   ----6
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 4)
	g.AddEdge(2, 5)
	g.AddEdge(2, 6)
	g.AddEdge(4, 5)
	g.AddEdge(4, 6)
	g.AddEdge(5, 6)
	ord, breaks := g.DegeneracyOrdering()
	fmt.Println("Degeneracy:", len(breaks)-1)
	fmt.Println("k-breaks:", breaks)
	fmt.Println("Ordering:", ord)
	for k, x := range breaks {
		fmt.Printf("nodes of %d-core(s): %d\n", k, ord[:x])
	}
	// Output:
	// Degeneracy: 3
	// k-breaks: [7 6 6 4]
	// Ordering: [4 5 6 2 0 1 3]
	// nodes of 0-core(s): [4 5 6 2 0 1 3]
	// nodes of 1-core(s): [4 5 6 2 0 1]
	// nodes of 2-core(s): [4 5 6 2 0 1]
	// nodes of 3-core(s): [4 5 6 2]
}

func ExampleUndirected_Degree() {
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

func ExampleUndirected_Density() {
	// 0---1
	// |
	// 2---3
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(2, 3)
	fmt.Println(g.Density())
	// Output:
	// 0.5
}

func ExampleUndirected_Eulerian_cycle() {
	//   0---
	//  /    \
	//  \     \
	//   1-----3
	//  / \   / \
	//  \  \ /  /
	//   ---2---
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 3)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	g.AddEdge(2, 3)
	fmt.Println(g.Eulerian())
	// Output:
	// -1 -1 <nil>
}

func ExampleUndirected_Eulerian_königsberg() {
	//   0--
	//  / \ \
	//  \ /  \
	//   1----3
	//  / \  /
	//  \ / /
	//   2--
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(0, 3)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	fmt.Println(g.Eulerian())
	// Output:
	// 0 1 non-Eulerian
}

func ExampleUndirected_Eulerian_loopIsolated() {
	//  0  1--\
	//      \-/
	var g graph.Undirected
	g.AddEdge(1, 1)
	fmt.Println(g.Eulerian())
	// Output:
	// -1 -1 <nil>
}

func ExampleUndirected_Eulerian_path() {
	//  /--\
	// 0----1
	//  \--/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 1)
	g.AddEdge(0, 1)
	fmt.Println(g.Eulerian())
	// Output:
	// 0 1 <nil>
}

func ExampleUndirected_EulerianCycle() {
	//   0---
	//  /    \
	//  \     \
	//   1-----3
	//  / \   / \
	//  \  \ /  /
	//   ---2---
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 3)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	g.AddEdge(2, 3)
	fmt.Println(g.EulerianCycle())
	// Output:
	// [0 1 3 2 1 2 3 0] <nil>
}

func ExampleUndirected_EulerianCycleD() {
	// 0----1
	//  \  /|\
	//   \ \|/
	//    --2--\
	//       \-/
	var g graph.Undirected
	// add 6 edges
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 2) // loop
	m := g.Size()
	fmt.Println("m =", m)
	fmt.Println(g.EulerianCycleD(m))
	// Output:
	// m = 6
	// [0 1 2 2 1 2 0] <nil>
}

func ExampleUndirected_EulerianPath() {
	//  /--\
	// 0----1
	//  \--/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 1)
	g.AddEdge(0, 1)
	fmt.Println(g.EulerianPath())
	// Output:
	// [0 1 0 1] <nil>
}

func ExampleUndirected_EulerianPathD() {
	//  /--\
	// 0----1
	//  \--/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 1)
	g.AddEdge(0, 1)
	fmt.Println(g.EulerianPathD(3, 0))
	// Output:
	// [0 1 0 1] <nil>
}

func ExampleUndirected_EulerianStart() {
	//   0--
	//  /   \
	//  \    \
	//   1----3
	//  / \  /
	//  \ / /
	//   2--
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(0, 3)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	fmt.Println(g.EulerianStart())
	// Output:
	// 2
}

func ExampleUndirected_IsConnected() {
	// undirected graph:
	//   0
	//  / \
	// 1   2
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// connected:  true
}

func ExampleUndirected_IsConnected_notConnected() {
	// undirected graph:
	//   0   1
	//  / \
	// 2   3
	var g graph.Undirected
	g.AddEdge(0, 2)
	g.AddEdge(0, 3)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// connected:  false
}
func ExampleUndirected_IsTree() {
	//  0--\
	//  |  |
	//  \--/   1   3
	//        /   / \
	//       2   4---5
	var g graph.Undirected
	g.AddEdge(0, 0)
	g.AddEdge(1, 2)
	g.AddEdge(3, 4)
	g.AddEdge(3, 5)
	g.AddEdge(4, 5)
	fmt.Println(g.IsTree(0))
	fmt.Println(g.IsTree(1))
	fmt.Println(g.IsTree(3))
	// Output:
	// false false
	// true false
	// false false
}

func ExampleUndirected_Size() {
	//   0--\
	//  / \-/
	// 1
	var g graph.Undirected
	g.AddEdge(0, 0)
	g.AddEdge(0, 1)
	fmt.Println("Size:", g.Size())
	fmt.Printf("(Arc size = %d)\n", g.ArcSize())
	// Output:
	// Size: 2
	// (Arc size = 3)
}
